package main

//This will be the main routing file for all internal operations
//such as those which don't create external services

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type viewState int
type focused int

const (
	mainView viewState = iota
	csvView
	createView
	observeView
	modView
	settingsView
	testView
)

const (
	fmvTable focused = iota
	clear
	vdList
	raceTable
)

type Tui struct {
	state   viewState
	focused focused
	help    help.Model

	createForm csvForm
	room       room

	list         list.Model
	vdSearch     list.Model
	vdSearchKeys vdSearchKeyMap

	//Components for assembling the Race Roster
	//colorGroups []table.Model
	fmvTable  table.Model
	fmvKeys   fmvTableKeyMap
	raceTable table.Model
	raceKeys  raceTableKeyMap

	raceRecord  raceRecord
	colorTables []table.Model // I know...

	velocidronePilots, registeredPilots []Pilot
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type mi string                 // front page menu link
func (mi) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(mi)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

func NewTui() *Tui {
	const defaultWidth = 20
	menuItems := []list.Item{
		mi("Create Race"),
		mi("Find Race"),
		mi("Moderate Race"),
		mi("Help / Settings"),
	}
	l := list.New(menuItems, itemDelegate{}, defaultWidth, 14)
	l.Title = "FMV RaceNite"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	sub := make(chan []byte)
	done := make(chan struct{})

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	conn := c

	room := room{
		sub:  sub,
		done: done,
		conn: conn,
	}

	return &Tui{list: l, state: mainView, fmvKeys: theFmvKeys, vdSearchKeys: theVdSearchKeys, raceKeys: theRaceTableKeys, help: help.New(), room: room}
}

func (m Tui) Init() tea.Cmd {
	return nil
}

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		}
		switch m.state {
		case mainView:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				switch m.list.SelectedItem().(mi) {
				case "Create Race":
					m.state = createView
					m, cmd := m.createForm.Update(formCmd())
					cmds = append(cmds, cmd, formCmd())
					return m, tea.Batch(cmds...)
				case "Find Race":
					findRoom, cmd := m.room.Update(msg)
					cmds = append(cmds, cmd, findRaceCmd())
					return findRoom, tea.Batch(cmds...)
				}
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case createView:
			switch m.focused {
			case fmvTable:
				switch {
				case key.Matches(msg, m.fmvKeys.checkin):
					if m.vdSearch.FilterState() != list.Filtering {
						x := m.fmvTable.SelectedRow()
						m.Checkin(x)
						fmvRows := updateFMVtable(m.registeredPilots)
						m.fmvTable.SetRows(fmvRows)
					}
				case key.Matches(msg, m.fmvKeys.checkinAll):
					if m.vdSearch.FilterState() != list.Filtering {
						for _, i := range m.registeredPilots {
							x := table.Row{i.DiscordName}
							m.CheckinAll(x)
						}
						fmvRows := updateFMVtable(m.registeredPilots)
						m.fmvTable.SetRows(fmvRows)
					}
				case key.Matches(msg, m.fmvKeys.switchToVd):
					m.fmvTable.Blur()
					m.focused = vdList
				}
			case vdList:
				switch {
				case key.Matches(msg, m.vdSearchKeys.addToFmv):
					if m.vdSearch.FilterState() != list.Filtering {
						listItem := m.vdSearch.SelectedItem().FilterValue()
						pilotFlag := false
						for _, i := range m.registeredPilots {
							if i.VdName == listItem {
								pilotFlag = true
								break
							}
						}
						if !pilotFlag {
							for _, i := range m.velocidronePilots {
								if i.VdName == listItem {
									var x Pilot
									x.DiscordName = i.VdName
									x.VdName = i.VdName
									x.QualifyingTime = i.QualifyingTime
									x.ModelName = i.ModelName
									m.registeredPilots = append(m.registeredPilots, x)
									fmvrows := updateFMVtable(m.registeredPilots)
									m.fmvTable.SetRows(fmvrows)
								}
							}
						}
					}
				case key.Matches(msg, m.vdSearchKeys.updateAtFmv):
					if m.vdSearch.FilterState() != list.Filtering {
						m.vdToFMVracer()
						fmvRows := updateFMVtable(m.registeredPilots)
						m.fmvTable.SetRows(fmvRows)
					}
				case key.Matches(msg, m.vdSearchKeys.switchToFmV):
					m.focused = fmvTable
					m.fmvTable.Focus()

				}

			}
			switch keypress := msg.String(); keypress {
			case "M", "m": // add error messege to view if too many pilots try to get pushed
				var counter = 0
				for _, pilot := range m.registeredPilots {
					if pilot.Status {
						counter++
					}
				}

				//////// LOGIC FOR SERVER PUT HERE ///////////
				if counter < 51 {
					m.raceRecord.Pilots = m.registeredPilots
					m.raceRecord.Round = 1
					m.raceRecord.Turn = 1

					/*
						m.raceTable = buildRaceTable()
						rows := updateRaceTable(m.raceRecord.Pilots)
						m.raceTable.SetRows(rows)
						m.state = modView
						m.focused = raceTable
						m.raceTable.Focus()

						sortedRacers := makeSortedRaceList(m.registeredPilots)
						groups := groupsArray(sortedRacers)
						m.colorTables = m.makeColorTables(groups)
						indexLen := len(groups)
						for i := 0; i < indexLen; i++ {
							rows := []table.Row{}
							for _, x := range groups[i] {
								rows = append(rows, x)
								m.colorTables[i].SetRows(rows)
							}
						}*/
					room, cmd := m.room.Update(msg)
					cmds = append(cmds, cmd, recordCmd(m.raceRecord)) // do not fuck with this
					return room, tea.Batch(cmds...)
				}
			}

			switch m.focused {
			case fmvTable:
				fmvstyle := fmvTableSelectedStyle("128", "207")
				m.fmvTable.SetStyles(fmvstyle)
				m.fmvTable, cmd = m.fmvTable.Update(msg)
				cmds = append(cmds, cmd)
				vdstyle := vdSearchSelectedStyle("242")
				m.vdSearch.SetDelegate(vdstyle)

				cmds = append(cmds, cmd)

			case vdList:
				style := fmvTableSelectedStyle("242", "249")
				m.fmvTable.SetStyles(style)
				vdstyle := vdSearchSelectedStyle("#EE6FF8")
				m.vdSearch.SetDelegate(vdstyle)
				//m.vdSearch, cmd = m.vdSearch.Update(msg)  <~~~~~ this breaks search function for VDlist
				//cmds = append(cmds, cmd)

			}
		case modView:
			switch {
			case key.Matches(msg, m.raceKeys.toMain):
				m.state = mainView
				m.focused = clear
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
			switch m.focused {
			case raceTable:
				m.raceTable, cmd = m.raceTable.Update(msg)
				cmds = append(cmds, cmd)
			}

		}
	case csvProcessedMsg: // vd, bound
		lists := msg
		m.registeredPilots = lists[1]
		m.fmvTable = buildFMVtable()
		rows := updateFMVtable(m.registeredPilots)
		m.fmvTable.SetRows(rows)
		m.velocidronePilots = lists[0]
		m.vdSearch = buildVelocidroneList(m.velocidronePilots)
		m.state = createView
		m.fmvTable, cmd = m.fmvTable.Update(msg)
		m.focused = fmvTable
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case testMsg:
		m.state = testView
	}

	if m.focused == vdList {
		m.vdSearch, cmd = m.vdSearch.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)

}

func (m Tui) View() string {
	if m.state != mainView {
		switch m.state {
		case createView:
			var footer string
			var titleColor string = "11"
			var checkedInPilots int

			for _, i := range m.registeredPilots {
				if i.Status {
					checkedInPilots++
				}
			}

			headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor))
			//paddingStyle := lipgloss.NewStyle().Padding(0, 2)

			fmvtitle := headerStyle.Render(fmt.Sprintf("\n\n FMV Voice Checkin (Pilots:%v)\n", checkedInPilots))
			fmvtableView := m.fmvTable.View()
			fmvBody := lipgloss.JoinVertical(lipgloss.Center, fmvtitle, fmvtableView)

			fmvBodyPadding := lipgloss.NewStyle().Padding(0, 14)
			vdSearchPadding := lipgloss.NewStyle().Padding(2, 0) //adjust the drop down
			vdSearchView := vdSearchPadding.Render(m.vdSearch.View())
			body := lipgloss.JoinHorizontal(lipgloss.Top, vdSearchView, fmvBodyPadding.Render(fmvBody), headerStyle.Render(fmvTag))

			switch m.focused {
			case vdList:
				footer = m.help.View(m.vdSearchKeys)

			case fmvTable:
				footer = m.help.View(m.fmvKeys)
			}

			footerPadding := lipgloss.NewStyle().PaddingLeft(2)
			viewPadding := lipgloss.NewStyle().PaddingLeft(6)
			view := lipgloss.JoinVertical(lipgloss.Left, body, footerPadding.Render(footer))
			return viewPadding.Render(view)

		case testView:
			return "Test View"
		case modView:
			colorNames := []string{"Gold", "Magenta", "Cyan", "Orange", "Green"}
			headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("207")).Padding(1, 0)
			header2Style := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Padding(1, 0, 0, 0).Underline(true)
			header := headerStyle.Render("FMV RaceNite Rawster")
			bodyPadding := lipgloss.NewStyle().Padding(0, 2)
			rtPadding := lipgloss.NewStyle().Padding(2, 0, 0, 0)

			rt := m.raceTable.View()
			raceTable := rtPadding.Render(lipgloss.JoinVertical(lipgloss.Center, header, rt))
			var groupTables []string
			for index, i := range m.colorTables {
				item := i.View()
				header := header2Style.Render(colorNames[index])
				table := lipgloss.JoinVertical(lipgloss.Center, header, item)
				groupTables = append(groupTables, table)

			}

			tables := lipgloss.JoinHorizontal(lipgloss.Center, groupTables...)
			footer := m.help.View(m.raceKeys)
			everything := bodyPadding.Render(lipgloss.JoinVertical(lipgloss.Left, raceTable, tables, footer))
			return everything
		case observeView:
			//Work on this next you dingbat
		}

	}

	body := m.list.View()
	return body
}

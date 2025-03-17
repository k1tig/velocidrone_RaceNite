package main

//This will be the main routing file for all internal operations
//such as those which don't create external services

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	vdList
	raceTable
)

type Tui struct {
	state   viewState
	focused focused

	createForm csvForm

	list     list.Model
	vdSearch list.Model

	//Components for assembling the Race Roster
	//colorGroups []table.Model
	fmvTable    table.Model
	raceTable   table.Model
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
		mi("Spectate"),
		mi("Moderate Race"),
		mi("Help / Settings"),
	}
	l := list.New(menuItems, itemDelegate{}, defaultWidth, 14)
	l.Title = "FMV RaceNite"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &Tui{list: l, state: mainView}
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
					m, cmd := m.createForm.Update(msg)
					cmds = append(cmds, cmd, formCmd())
					return m, tea.Batch(cmds...)
				}
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case createView:
			switch keypress := msg.String(); keypress {
			case "tab":
				if m.focused == fmvTable {
					m.fmvTable.Blur()
					m.focused = vdList
				} else {
					m.focused = fmvTable
					m.fmvTable.Focus()
				}
			case "a", "A": // adds VdPilot to FMVbound list
				if m.focused == vdList {
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
				}
			case "E", "e": // updates FMV users info from VD list
				if m.focused == vdList {
					if m.vdSearch.FilterState() != list.Filtering {
						m.vdToFMVracer()
						fmvRows := updateFMVtable(m.registeredPilots)
						m.fmvTable.SetRows(fmvRows)
					}
				}
			case "C":
				if m.focused == fmvTable {
					if m.vdSearch.FilterState() != list.Filtering {
						for _, i := range m.registeredPilots {
							x := table.Row{i.DiscordName}
							m.CheckinAll(x)
						}
						fmvRows := updateFMVtable(m.registeredPilots)
						m.fmvTable.SetRows(fmvRows)
					}
				}
			case "c":
				if m.focused == fmvTable {
					if m.vdSearch.FilterState() != list.Filtering {
						x := m.fmvTable.SelectedRow()
						m.Checkin(x)
						fmvRows := updateFMVtable(m.registeredPilots)
						m.fmvTable.SetRows(fmvRows)
					}
				}
			case "M", "m":
				m.raceTable = buildRaceTable()
				rows := updateRaceTable(m.registeredPilots)
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
				}

				for i := 0; i < indexLen-1; i++ {
					m.colorTables[i].Blur()
				}

			}
			switch m.focused {
			case fmvTable:
				m.fmvTable, cmd = m.fmvTable.Update(msg)
				cmds = append(cmds, cmd)
			case vdList:
				//m.vdSearch, cmd = m.vdSearch.Update(msg)  <~~~~~ this breaks search function for VDlist
				//cmds = append(cmds, cmd)

			}
		case modView:
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
	var helpText string
	if m.state != mainView {
		switch m.state {
		case createView:

			switch m.focused {
			case vdList:
				helpText = "~HELP KEYS~\nTab: focus table, A/a: add to fmv, E/e: assign vd item to FMV racer"
			}

			header := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
			padding := lipgloss.NewStyle().Padding(0, 2)
			listpadding := lipgloss.NewStyle().Padding(2, 6)
			num := "placeholder num"

			fmvtitle := header.Render(fmt.Sprintf("\n\n FMV Voice Checkin (count:%s)\n", num))
			fmvTable := m.fmvTable.View()
			fmvBody := padding.Render(lipgloss.JoinVertical(lipgloss.Center, fmvtitle, fmvTable))

			vdList := listpadding.Render(m.vdSearch.View())
			body := lipgloss.JoinHorizontal(lipgloss.Top, vdList, fmvBody, fmvTag)
			view := lipgloss.JoinVertical(lipgloss.Left, body, listpadding.Render(helpText))
			return view

		case testView:
			return "Test View"
		case modView:

			headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("207")).Padding(2, 0)
			header := headerStyle.Render("FMV RaceNite Rawster")
			padding := lipgloss.NewStyle().Padding(1, 6)

			body := padding.Render(m.raceTable.View())
			view := lipgloss.JoinVertical(lipgloss.Center, header, body)
			var groupTables []string
			for _, i := range m.colorTables {
				view := i.View()
				groupTables = append(groupTables, view)

			}
			tables := lipgloss.JoinHorizontal(lipgloss.Center, groupTables...)
			everything := lipgloss.JoinVertical(lipgloss.Center, view, tables)
			return everything
		}

	}
	return "\n" + m.list.View()
}

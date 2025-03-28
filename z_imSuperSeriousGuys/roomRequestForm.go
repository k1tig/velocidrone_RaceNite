package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type roomstate int

// msgs
type recordMsg struct {
	raceRecord raceRecord
}

func recordCmd(rr raceRecord) tea.Cmd {
	return func() tea.Msg {
		return recordMsg{raceRecord: rr}
	}
}

type initRaceTableMsg struct {
	table table.Model
}

func initRaceTableCmd(pilots []Pilot) tea.Cmd {
	initTable := buildRaceTable()
	rows := updateRaceTable(pilots)
	initTable.SetRows(rows)

	return func() tea.Msg {
		return initRaceTableMsg{table: initTable}
	}
}

type findRaceMsg struct{}

func findRaceCmd() tea.Cmd {
	return func() tea.Msg {
		return findRaceMsg{}
	}
}

const (
	defaultstate roomstate = iota
	formstate
	viewstate
)

var (
	rooms = []string{"1", "3", "5", "7"}
	Dingy = "yes"
)

type room struct {
	state    roomstate
	help     help.Model
	raceKeys raceTableKeyMap

	form   *huh.Form
	roomId int

	//roomKey     int
	raceTable   table.Model
	colorTables []table.Model
	raceRecord  raceRecord
}

func newRoomForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("raceid").
				Title("Active Races:").
				OptionsFunc(func() []huh.Option[string] {
					s := rooms
					// simulate API call
					time.Sleep(1 * time.Second)
					return huh.NewOptions(s...)
				},
					huh.NewConfirm().
						Key("done").
						Title("All done?").
						Validate(func(v bool) error {
							if !v {
								return fmt.Errorf("finish it")
							}
							return nil
						}).
						Affirmative("Yep").
						Negative("Wait, no"),
				),
		).WithWidth(45).
			WithShowHelp(false).
			WithShowErrors(false),
	)
	return form
}
func (m room) Init() tea.Cmd {
	return nil
}
func (m room) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	case findRaceMsg:
		m.form = newRoomForm()
		m.state = formstate
		//return m, tea.Batch(cmds...)

	case recordMsg:
		m.raceRecord = msg.raceRecord
		return m, initRaceTableCmd(m.raceRecord.Pilots)
	case initRaceTableMsg:
		m.raceTable = msg.table
		sortedPilots := makeSortedRaceList(m.raceRecord.Pilots)
		pilotGroups := groupsArray(sortedPilots)

		m.colorTables = makeColorTables(pilotGroups)
		indexLen := len(pilotGroups)
		for i := 0; i < indexLen; i++ {
			rows := []table.Row{}
			for _, x := range pilotGroups[i] {
				rows = append(rows, x)
				m.colorTables[i].SetRows(rows)
			}
			m.state = viewstate

			m.raceTable, cmd = m.raceTable.Update(msg)

			cmds = append(cmds, cmd)

		}
		return m, tea.Batch(cmds...)
	}

	switch m.state {

	case formstate:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}
		if m.form.State == huh.StateCompleted {
			if m.state == formstate {
				x := m.form.GetString("raceid")
				//// seriously don't forget raceKey
				id, err := strconv.Atoi(x)
				if err != nil {
					fmt.Println("Error during conversion:", err)
					return m, nil
				}
				m.roomId = id
				//tea.Msg to go get raceRecord from server
			}
			// **don't forget to add option for raceKey (modView)
			// get race record to build race tables with + color tables
			//
		}

	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
func (m room) View() string {
	switch m.state {
	case defaultstate:
		return "default state"
	case formstate:
		return m.form.View()

	case viewstate:
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
	}

	return fmt.Sprint(m.state, Dingy)
}

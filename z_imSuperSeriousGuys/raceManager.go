package main

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type roomstate int
type recordMsg struct{ raceRecord raceRecord }
type initRaceTableMsg struct{ table table.Model }
type findRaceMsg struct{}
type Message struct {
	Event      string          `json:"event"`
	Parameters json.RawMessage `json:"parameters"`
}

const (
	defaultstate roomstate = iota
	formstate
	viewstate
	teststate
)

var (
	rooms = []string{"1", "3", "5", "7"}
	Dingy = "yes"
)

type room struct {
	state    roomstate
	help     help.Model
	raceKeys raceTableKeyMap
	form     *huh.Form
	conn     *websocket.Conn
	sub      chan []byte
	done     chan struct{}
	//roomKey     int // permision to mod racesRecords on server

	raceTable   table.Model
	colorTables []table.Model
	raceRecord  raceRecord
	testMsg     []byte
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
			closeWebSocket(m.conn)
			return m, tea.Quit
		}
	//msg cases
	case findRaceMsg:
		rooms = getRaceRecords()
		m.form = newRoomForm()
		m.state = formstate

	case recordMsg:
		m.raceRecord = msg.raceRecord
		sendRaceRecord(m.raceRecord)
		cmds = append(cmds, m.initWsReader(), waitForMsg(m.sub), initRaceTableCmd(m.raceRecord.Pilots))
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
		//might be broke as shit now
	case responseMsg:
		var rr raceRecord
		err := json.Unmarshal(msg, &rr)
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		m.raceRecord = rr
		cmds = append(cmds, initRaceTableCmd(m.raceRecord.Pilots), waitForMsg(m.sub))
	}

	//state switches
	switch m.state {
	case formstate:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		///////////fix this shit
		if m.form.State == huh.StateCompleted {
			x := m.form.GetString("raceid")
			m.raceRecord = getRaceRecordsById(x)
			cmds = append(cmds, m.initWsReader(), waitForMsg(m.sub), initRaceTableCmd(m.raceRecord.Pilots))
			return m, tea.Batch(cmds...)
		}

	case viewstate:
		m.raceTable, cmd = m.raceTable.Update(msg)
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
func (m room) View() string {
	switch m.state {
	case teststate:
		message := "Recieved: " + string(m.testMsg)
		return message
	case defaultstate:
		return "Connecting...."
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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
	return tea.Batch(m.initWsReader(), waitForMsg(m.sub))
}
func (m room) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	//msg cases
	case findRaceMsg:
		rooms = getRaceRecords()
		m.form = newRoomForm()
		m.state = formstate
	case recordMsg:
		m.raceRecord = msg.raceRecord
		m.sendRaceRecord()
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

		//might be broke as shit now
	case roomMsg:
		m.sub = msg.room.sub
		m.done = msg.room.done
		m.conn = msg.room.conn
		return m, tea.Batch(m.initWsReader(), waitForMsg(m.sub))
	case responseMsg:
		var message Message
		err := json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		switch message.Event {
		case "update":
			var rr raceRecord
			err := json.Unmarshal(message.Parameters, &rr)
			if err != nil {
				fmt.Printf("err: %s", err)
			}
			m.raceRecord = rr
			return m, tea.Batch(initRaceTableCmd(rr.Pilots), m.initWsReader(), waitForMsg(m.sub))
		}
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
			return m, tea.Batch(initRaceTableCmd(m.raceRecord.Pilots)) ///////thisssss is where its fuckeeedd
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

func (m room) sendRaceRecord() {
	url := "http://localhost:8080/brackets"
	post := m.raceRecord

	requestBody, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return
	}

}

func getRaceRecords() []string {
	type records []raceRecord
	var recordsData records
	var raceIds []string
	url := "http://localhost:8080/brackets"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	err = json.Unmarshal(body, &recordsData)
	if err != nil {
		fmt.Println("error unmarshalling json from server:", err)
	}

	for _, raceId := range recordsData {
		raceIds = append(raceIds, strconv.Itoa(raceId.Id))
	}
	return raceIds
}

func getRaceRecordsById(id string) raceRecord {
	var record raceRecord
	url := "http://localhost:8080/brackets/" + id
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	err = json.Unmarshal(body, &record)
	if err != nil {
		fmt.Println("error unmarshalling json from server:", err)
	}
	return record
}

func recordCmd(rr raceRecord) tea.Cmd {
	return func() tea.Msg {
		return recordMsg{raceRecord: rr}
	}
}
func initRaceTableCmd(pilots []Pilot) tea.Cmd {
	initTable := buildRaceTable()
	rows := updateRaceTable(pilots)
	initTable.SetRows(rows)
	return func() tea.Msg { return initRaceTableMsg{table: initTable} }
}
func findRaceCmd() tea.Cmd {
	return func() tea.Msg { return findRaceMsg{} }
}

type roomMsg struct{ room room }

/*
func initWs() tea.Cmd {
	return func() tea.Msg {
		sub := make(chan []byte)
		done := make(chan struct{})
		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
		log.Printf("connecting to %s", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		conn := c

		var room = room{
			sub:  sub,
			done: done,
			conn: conn,
		}

		return roomMsg{
			room: room,
		}
	}

}*/

// /////////////////////////////////////////////////////////////   this is broken /////////////////////////////////
func (m room) initWsReader() tea.Cmd {
	// Start listening for WebSocket messages in a goroutine
	//return func() tea.Msg {

	return func() tea.Msg {
		for {
			_, message, err := m.conn.ReadMessage()
			if err != nil {
				log.Fatal("init reader err")
			}
			m.sub <- message
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type responseMsg []byte

func waitForMsg(sub chan []byte) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

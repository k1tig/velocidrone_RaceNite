package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

// tea Msgs and Cmds
type csvProcessedMsg [][]Pilot
type testMsg struct{}

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

func (m room) initWsReader() tea.Cmd {
	return func() tea.Msg {
		defer close(m.done)
		for {
			_, message, err := m.conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			m.sub <- message
		}
	}
}

type responseMsg []byte

func waitForMsg(sub chan []byte) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

//////// end msg cmds /////

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
func closeWebSocket(conn *websocket.Conn) {
	err := conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second))
	if err != nil && err != websocket.ErrCloseSent {
		log.Println("write close error:", err)
		return
	}
	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return
	}

	for {
		_, _, err = conn.NextReader()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break
		}
		if err != nil {
			break
		}
	}
	err = conn.Close()
	if err != nil {
		log.Println("close error:", err)
		return
	}
}

func sendRaceRecord(record raceRecord) {
	url := "http://localhost:8080/brackets"
	post := record

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

func fmvTableSelectedStyle(bgColor, fgColor string) table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("ffb3fd")).
		Foreground(lipgloss.Color("239")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor))

	return s
}

func vdSearchSelectedStyle(color string) list.DefaultDelegate {
	s := list.NewDefaultDelegate()
	s.Styles.SelectedTitle = s.Styles.SelectedTitle.
		Foreground(lipgloss.Color(color)).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: color})
	s.Styles.SelectedDesc = s.Styles.SelectedTitle

	return s
}
func buildRaceTable() table.Model {
	fmvColumns := []table.Column{
		{Title: "Pilot", Width: 16},
		{Title: "Points", Width: 7},
		{Title: "Qualify Time", Width: 12},
		{Title: "H1", Width: 8},
		{Title: "H2", Width: 8},
		{Title: "H3", Width: 8},
		{Title: "H4", Width: 8},
		{Title: "H5", Width: 8},
		{Title: "H6", Width: 8},
		{Title: "H7", Width: 8},
		{Title: "H8", Width: 8},
		{Title: "H9", Width: 8},
		{Title: "H10", Width: 8},
	}
	rows := []table.Row{}
	raceTable := table.New(
		table.WithColumns(fmvColumns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(16),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("ffb3fd")).
		Foreground(lipgloss.Color("239")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color("128"))

	raceTable.SetStyles(s)

	return raceTable
}

func buildFMVtable() table.Model {
	fmvColumns := []table.Column{
		{Title: "Pilot", Width: 16},
		{Title: "VD Name", Width: 16},
		{Title: "Qualify time", Width: 16},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{}

	fmvTable := table.New(
		table.WithColumns(fmvColumns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("ffb3fd")).
		Foreground(lipgloss.Color("239")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color("128")).
		Foreground(lipgloss.Color("207"))

	fmvTable.SetStyles(s)

	return fmvTable
}
func updateRaceTable(pilots []Pilot) []table.Row {
	rows := []table.Row{}

	//Maybe this will work?
	sort.Slice(pilots, func(i, j int) bool {
		var a, b float64
		a, _ = strconv.ParseFloat(pilots[i].QualifyingTime, 64)
		b, _ = strconv.ParseFloat(pilots[j].QualifyingTime, 64)
		return a < b
	})

	for _, i := range pilots {
		if i.Status {
			var s []string
			var fakePoints string
			vdName := i.VdName
			fakePoints = "0"
			qTime := i.QualifyingTime
			rt := floatListToStringList(i.RaceTimes)
			s = append(s, vdName, fakePoints, qTime)
			s = append(s, rt...)
			rows = append(rows, s)
		}
	}

	return rows

}

func updateFMVtable(racers []Pilot) []table.Row {
	rows := []table.Row{}

	for _, i := range racers {
		var s []string
		var status string
		name := i.DiscordName
		vdName := i.VdName
		qtime := i.QualifyingTime
		if !i.Status {
			status = "-"
		} else {
			status = "Entered"
		}
		s = append(s, name, vdName, qtime, status)
		rows = append(rows, s)
	}
	return rows
}

func buildVelocidroneList(vdSheet []Pilot) list.Model {
	var racers = []list.Item{}
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("242")).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "242"})
	d.Styles.SelectedDesc = d.Styles.SelectedTitle
	vdList := list.New(racers, d, 0, 0)
	vdList.Title = "Velocidrone Sheet"
	vdList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		UnsetBackground()
	//Bold(true).
	//Underline(true)

	vdList.SetSize(28, 20)

	for _, racer := range vdSheet {
		obj := listRacer{name: racer.VdName, time: racer.QualifyingTime, craft: racer.ModelName}
		//items = append(items, obj)
		vdList.InsertItem(99999, obj) //out of range placement appends item to list
	}
	return vdList
}

func processForm(e *huh.Form) (vd, fmvBound []Pilot) {
	discordTarget := GetDiscordId(e.GetString("discord"))
	fmvTarget := GetFMVvoice(e.GetString("fmv"))
	vdTarget := GetVdRacers(e.GetString("vd"))
	registeredTarget := BindLists(vdTarget, fmvTarget, discordTarget)
	return vdTarget, registeredTarget

}

func makeSortedRaceList(pilotList []Pilot) [][]string {
	type cleanRacer struct {
		racer  string
		time   float64
		points float64
	}
	var (
		cleanRacers []cleanRacer
		racers      []Pilot
		racingList  [][]string
	)

	for _, pilot := range pilotList {
		if pilot.Status {
			racers = append(racers, pilot)
		}
	}

	for _, pilot := range racers {
		var cr cleanRacer
		cr.racer = pilot.VdName
		cr.points = pilot.Points
		QualifyingTime, _ := strconv.ParseFloat(pilot.QualifyingTime, 64)
		cr.time = QualifyingTime

		cleanRacers = append(cleanRacers, cr)
	}

	sort.Slice(cleanRacers, func(i, j int) bool {
		return cleanRacers[i].time < cleanRacers[j].time
	})

	for _, i := range cleanRacers {
		var racestring []string
		racestring = append(racestring, i.racer)
		raceTime := strconv.FormatFloat(i.time, 'g', 3, 64)
		points := strconv.FormatFloat(i.points, 'g', 3, 64)

		racestring = append(racestring, raceTime)
		racestring = append(racestring, points)
		racingList = append(racingList, racestring)
	}
	return racingList
}

func groupsArray(vdList [][]string) [][][]string {
	//makes a group of groups with the total amount of racers not exceeding a +1 differential
	var maxGroupsize = 8
	var grouplength int
	var totalGroups int
	var modulus int
	var racers = (len(vdList))

	if racers > 40 {
		maxGroupsize = 10
	}

	for i := 1; i <= maxGroupsize; i++ {
		if float64(racers)/float64(i) <= float64(maxGroupsize) { //  42_1_2_3_4_5....oh its a float rounding issue...moron. note:fixed*
			totalGroups = i
			modulus = int(racers) % int(i)
			if modulus == 0 {
				grouplength = racers / i
			} else {
				grouplength = (racers - modulus) / i
			}
			break
		}
	}

	var groupStructure = make([][][]string, int(totalGroups))
	var c int
	x := modulus

	for i := 1; i <= totalGroups; i++ {

		if x > 0 { // distribues the modulus between the lower teir groups
			racers := vdList[c : int(i)*(int(grouplength)+1)]
			groupStructure[int(i)-1] = racers
			x--
			c += int(grouplength) + 1
		} else { // groups that don't take a modulus
			racers := vdList[c : c+int(grouplength)]
			groupStructure[int(i)-1] = racers
			c += int(grouplength)
		}
	}
	return groupStructure
}

func makeColorTables(brackets [][][]string) (tableList []table.Model) {
	indexLen := len(brackets)
	colors := []string{"190", "171", "123", "214", "47"}

	for i := 0; i < indexLen; i++ {
		columns := []table.Column{
			{Title: "Pilot", Width: 16},
			{Title: "Time", Width: 7},
			{Title: "Points", Width: 7},
		}
		rows := []table.Row{}
		groupTable := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithHeight(14),
		)
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("ffb3fd")).
			Foreground(lipgloss.Color("239")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Background(nil).Bold(false)

		s.Cell = s.Cell.
			Foreground(lipgloss.Color(colors[i]))
			//Foreground(lipgloss.Color("128"))

		groupTable.SetStyles(s)
		tableList = append(tableList, groupTable)
	}
	return tableList
}

func floatListToStringList(floatList [10]float64) []string {
	stringList := make([]string, len(floatList))
	for i, num := range floatList {
		if num == 0 {

		}
		stringList[i] = strconv.FormatFloat(num, 'f', 3, 64)
	}
	for i, str := range stringList {
		if str == "0.000" {
			stringList[i] = "-"
		}

	}
	return stringList
}

// model methods

func (m Tui) vdToFMVracer() {
	r := m.fmvTable.SelectedRow()
	listItem := m.vdSearch.SelectedItem().FilterValue()
	for index, i := range m.registeredPilots {
		if r[0] == i.DiscordName {
			for _, x := range m.velocidronePilots {
				if x.VdName == listItem {
					i.VdName = x.VdName
					i.QualifyingTime = x.QualifyingTime
					i.ModelName = x.ModelName
					m.registeredPilots[index] = i
					return
				}
			}
		}
	}
}

func (m Tui) Checkin(r table.Row) {
	for index, i := range m.registeredPilots {
		if r[0] == i.DiscordName {
			switch i.Status {
			case !true:
				if i.QualifyingTime != "CHECK IN Please!" {
					i.Status = true
				}
			case true:
				i.Status = false
			}
			m.registeredPilots[index] = i
		}
	}
}

func (m Tui) CheckinAll(r table.Row) {
	for index, i := range m.registeredPilots {
		if r[0] == i.DiscordName {
			if i.QualifyingTime != "CHECK IN Please!" {
				i.Status = true
			}
		}
		m.registeredPilots[index] = i
	}
}

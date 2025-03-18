package main

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type csvProcessedMsg [][]Pilot

func processForm(e *huh.Form) (vd, fmvBound []Pilot) {
	discordTarget := GetDiscordId(e.GetString("discord"))
	fmvTarget := GetFMVvoice(e.GetString("fmv"))
	vdTarget := GetVdRacers(e.GetString("vd"))
	registeredTarget := BindLists(vdTarget, fmvTarget, discordTarget)
	return vdTarget, registeredTarget

}

func buildVelocidroneList(vdSheet []Pilot) list.Model {
	var racers = []list.Item{}
	vdList := list.New(racers, list.NewDefaultDelegate(), 0, 0)
	vdList.Title = "Velocidrone Sheet"
	vdList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("178")).
		Background(lipgloss.Color("0")).
		Bold(true).
		Underline(true)
	vdList.SetSize(28, 20)
	for _, racer := range vdSheet {
		obj := listRacer{name: racer.VdName, time: racer.QualifyingTime, craft: racer.ModelName}
		//items = append(items, obj)
		vdList.InsertItem(99999, obj) //out of range placement appends item to list
	}
	return vdList
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
		Background(lipgloss.Color("128"))
	fmvTable.SetStyles(s)

	return fmvTable
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

type testMsg struct{}

/*func testCmd() tea.Msg {
	return testMsg{}
}*/

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
func makeSortedRaceList(pilotList []Pilot) [][]string {
	type cleanRacer struct {
		racer  string
		time   float64
		points float64
	}

	var cleanRacers []cleanRacer
	var racers []Pilot

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

	var racingList [][]string
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

	racers := len(vdList)
	if racers > 40 {
		maxGroupsize = 10
	}

	for i := 1; i <= maxGroupsize; i++ {
		if racers/i <= maxGroupsize {
			totalGroups = i
			modulus = racers % i
			if modulus == 0 {
				grouplength = racers / i
			} else {
				grouplength = (racers - modulus) / i
			}
			break
		}
	}

	var groupStructure = make([][][]string, totalGroups)
	var c int
	x := modulus

	for i := 1; i <= totalGroups; i++ {
		if x > 0 { // distribues the modulus between the lower teir groups
			racers := vdList[c : i*(grouplength+1)]
			groupStructure[i-1] = racers
			x--
			c += grouplength + 1
		} else { // groups that don't take a modulus
			racers := vdList[c : c+grouplength]
			groupStructure[i-1] = racers
			c += grouplength
		}
	}
	return groupStructure
}

func (m Tui) makeColorTables(brackets [][][]string) (tableList []table.Model) {
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

package main

import (
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
		table.WithHeight(20),
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

func updateRaceTable(racers []Pilot) []table.Row {
	rows := []table.Row{}

	for _, i := range racers {
		if i.Status {
			var s []string
			vdName := i.VdName
			qTime := i.QualifyingTime
			rt := floatListToStringList(i.RaceTimes)

			s = append(s, vdName, qTime)
			s = append(s, rt...)
			rows = append(rows, s)
		}
	}
	return rows

}
func floatListToStringList(floatList [10]float64) []string {
	stringList := make([]string, len(floatList))
	for i, num := range floatList {
		stringList[i] = strconv.FormatFloat(num, 'f', 3, 64)
	}
	return stringList
}

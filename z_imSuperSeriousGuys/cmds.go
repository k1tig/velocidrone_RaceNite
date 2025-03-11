package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type csvProcessedMsg [][]Pilot

func processForm(e *huh.Form) (d1, f1, v1, reg1 []Pilot) {
	discordTarget := GetDiscordId(e.GetString("discord"))
	fmvTarget := GetFMVvoice(e.GetString("fmv"))
	vdTarget := GetVdRacers(e.GetString("vd"))
	registeredTarget := BindLists(vdTarget, fmvTarget, discordTarget)
	return discordTarget, fmvTarget, vdTarget, registeredTarget

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
	vdList.SetSize(20, 20)
	for _, racer := range vdSheet {
		obj := listRacer{name: racer.VdName, time: racer.QualifyingTime, craft: racer.ModelName}
		//items = append(items, obj)
		vdList.InsertItem(99999, obj) //out of range placement appends item to list
	}
	return vdList
}

func buildFMVtable(boundList []Pilot) table.Model {
	fmvColumns := []table.Column{
		{Title: "Pilot", Width: 16},
		{Title: "VD Name", Width: 16},
		{Title: "Qualify time", Width: 16},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{}

	for _, i := range boundList {
		var s []string
		var status string
		name := i.DiscordName
		vdName := i.VdName
		qtime := i.QualifyingTime
		if i.Status == NulPilot.Status {
			status = "-"
		} else {
			status = i.Status
		}
		s = append(s, name, vdName, qtime, status)
		rows = append(rows, s)
	}

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

type testMsg struct{}

/*func testCmd() tea.Msg {
	return testMsg{}
}*/

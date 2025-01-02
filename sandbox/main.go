package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type racer struct {
	name            string
	r1, r2, r3      string
	r4, r5, r6      string
	r7, r8, r9, r10 string
	s1, s2, s3      string
	p1, p2, p3      string
	p4, p5, p6      string
	p7, p8, p9, p10 string
	B1, B2, B3      string
}

type model struct {
	Split_1 table.Model
	racers  []racer
}

func initModel() model {
	columns := []table.Column{
		{Title: "Racer", Width: 10},
		{Title: "Bracket Time", Width: 14},
		{Title: "H1 Time", Width: 8},
		{Title: "H1 Pts", Width: 8},
		{Title: "H2 Time", Width: 8},
		{Title: "H2 Pts", Width: 8},
		{Title: "H3 Time", Width: 8},
		{Title: "H3 Pts", Width: 8},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(10),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("3")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("220")).
		Background(lipgloss.Color("236")).
		Bold(false)
	t.SetStyles(s)

	r := []string{"asiy", "MOTODRONEX", "eedok", "uGeLLin", "andyy", "MGescapades", "RoflCopter!", "AP3X",
		"Not Sure", "Barnyard", "Mayan_Hawk", ".MrE.", "MrMan", "XaeroFPV", "Kuzyatron", "Zikefire",
		"jon E5", "DeMic", "DreadPool", "Lounds", "PilotInCommand", "SilkJamFPV", "TEDDY_TUNED", "Runnin_Lizzard",
		"Pweeen", "landsquid", "kiz", "CAVEMAN_", "boondockstryker", "k1itg", "Derpy Hooves", "Charlito",
		"Sillybutter", "waffle3_0", "Dogbowl", "WalkerFPV", "BallHawk", "Hyde", "timmah1991", "FPVMartin",
	}
	var rList []racer
	for _, x := range r {
		var i racer
		i.name = x
		rList = append(rList, i)
	}

	m := model{t, rList}

	return m
}
func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Split_1.Focused() {
				m.Split_1.Blur()
			} else {
				m.Split_1.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.Split_1.SelectedRow()[1]),
			)
		case "t":
			x := m.buildTable()
			m.Split_1.SetRows(x)

		case "r":
			var e int
			for i, x := range m.racers {
				if x.name == "eedok" {
					e = i
				}
			}
			m.racers[e].r1 = "69"
			x := m.buildTable()
			m.Split_1.SetRows(x)

		}

	}
	m.Split_1, cmd = m.Split_1.Update(msg)
	return m, cmd
}
func (m model) View() string {

	return baseStyle.Render(m.Split_1.View()) + "\n"
}

func main() {
	if _, err := tea.NewProgram(initModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) buildTable() []table.Row {
	newRows := []table.Row{}
	for _, i := range m.racers {
		x := []string{i.name, i.B1, i.r1, i.p1, i.r2, i.p2, i.r3, i.p3}
		newRows = append(newRows, x)

	}

	return newRows

}

//Function to update racers from a list
/*
func (m model) qualifying() {
	type vd struct {
		name string
		time string
	}
	var vdList []vd
	for e, i := range m.racers {
		for _, x := range vdList {
			if i.name == x.name {
				m.racers[e].r1 = x.time
			}
		}

	}
}
*/

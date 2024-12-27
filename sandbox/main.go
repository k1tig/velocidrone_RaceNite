package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table []table.Model
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		rows := addRows()
		m.table[1].SetRows(rows[0])
		m.table[0].SetRows(rows[1])
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table[1].Focused() {
				m.table[1].Blur()
			} else {
				m.table[1].Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table[1].SelectedRow()[1]),
			)
		}
	}
	m.table[1], cmd = m.table[1].Update(msg)
	return m, cmd
}

func (m model) View() string {
	leftText := "Magenta Group"
	left := baseStyle.Render(m.table[1].View())
	left = lipgloss.JoinVertical(lipgloss.Center, leftText, left)
	rightText := "Gold Group"
	right := baseStyle.Render(m.table[0].View())
	right = lipgloss.JoinVertical(lipgloss.Center, rightText, right)
	body := lipgloss.JoinHorizontal(lipgloss.Left, left, right)
	return body
}

func main() {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}
	rows := []table.Row{}
	t1 := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	t2 := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t1.SetStyles(s)
	t2.SetStyles(s)
	type tables []table.Model
	m := model{tables{t1, t2}}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func addRows() [][]table.Row {
	var listOfRows [][]table.Row
	rows1 := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
		{"7", "Cairo", "Egypt", "21,750,020"},
	}
	rows2 := []table.Row{
		{"8", "a", "Japan", "37,274,000"},
		{"9", "b", "India", "32,065,760"},
		{"10", "c", "China", "28,516,904"},
		{"11", "D", "Bangladesh", "22,478,116"},
		{"12", "ã Paulo", "Brazil", "22,429,800"},
		{"13", "M", "Mexico", "22,085,140"},
		{"14", "z", "Egypt", "21,750,020"},
	}
	listOfRows = append(listOfRows, rows1, rows2)
	return listOfRows
}

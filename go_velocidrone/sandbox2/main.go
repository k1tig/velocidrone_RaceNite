package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

type state uint

const (
	listView state = iota
	tableView
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var titleStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("288"))

type model struct {
	table table.Model
	list  list.Model
	state state
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		case "tab":
			if m.state == tableView {
				m.table.Blur()
				m.state = listView

			} else {
				m.state = tableView
				m.table.Focus()
				return m, cmd
			}
		}
		switch m.state {
		case listView:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case tableView:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	T1 := "\n\n\n\n\n" + titleStyle.Render("FMV RaceNite Racers Entered") + "\n" + baseStyle.Render(m.table.View()) + "\n\n\n" + m.table.HelpView() + "\n"
	L1 := docStyle.Render(m.list.View())
	base := lipgloss.JoinHorizontal(lipgloss.Top, L1, T1)
	return "\n\n\n" + base
}

func main() {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 10},
		{Title: "Craft", Width: 10},
		{Title: "Time", Width: 10},
	}

	rows := []table.Row{
		{"1", "Blasta", "TBS Spec", "52.33"},
		{"2", "Eedok", "TBS Spec", "52.34"},
		{"3", "Mr E", "Twig XL", "53.11"},
		{"5", "Xaero", "TBS Spec", "54.49"},
		{"6", "Jon E5", "Twig XL", "57.44"},
		{"9", "Blunty", "Twig XL", "60.18"},
		{"21", "Sillybutter", "Twig XL", "69.49"},
		{"25", "TreeSeeker", "Twig XL", "90.12"},
		{"55", "DeMic", "Pinetree", "420.69"},
		{"56", "Blasta", "TBS Spec", "52.33"},
		{"58", "Eedok", "TBS Spec", "52.34"},
		{"60", "Mr E", "Twig XL", "53.11"},
		{"61", "Xaero", "TBS Spec", "54.49"},
		{"66", "Jon E5", "Twig XL", "57.44"},
		{"97", "Blunty", "Twig XL", "60.18"},
		{"251", "Sillybutter", "Twig XL", "69.49"},
		{"2455", "TreeSeeker", "Twig XL", "90.12"},
		{"5544", "DeMic", "Pinetree", "420.69"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(14),
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
	t.SetStyles(s)

	items := []list.Item{
		item{title: "Mr Man", desc: "69.33 - TBS Spec"},
		item{title: "SlikJam", desc: "44.10 - XL Twig"},
		item{title: "Bitter melon", desc: "102.33 - XL Twig"},
		item{title: "Mr Man", desc: "69.33 - TBS Spec"},
		item{title: "SlikJam", desc: "44.10 - XL Twig"},
		item{title: "Bitter melon", desc: "102.33 - XL Twig"},
		item{title: "Mr Man", desc: "69.33 - TBS Spec"},
		item{title: "SlikJam", desc: "44.10 - XL Twig"},
		item{title: "Bitter melon", desc: "102.33 - XL Twig"},
	}

	x := list.New(items, list.NewDefaultDelegate(), 0, 0)
	x.DisableQuitKeybindings()
	x.Title = "Select FMV Racers for FMV Roster"

	m := model{table: t,
		list: x,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

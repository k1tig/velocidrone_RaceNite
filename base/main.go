package main

import (
	"fmt"
	"os"

	rt "abc.com/base/racetools"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type vdracer struct {
	name, qTime, craft string
}
type fmvracer struct {
	name, qTime, craft string
}
type state uint

type model struct {
	velocidrone list.Model
	fmv         list.Model
	state       state
}

const (
	vdView state = iota
	fmvView
)

var (
	docStyle  = lipgloss.NewStyle().Margin(1, 2)
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("178")).
			Bold(true).
			Underline(true)
)
var check = "✓"

func (i vdracer) Title() string       { return i.name }
func (i vdracer) Description() string { return i.qTime + " | " + i.craft }
func (i vdracer) FilterValue() string { return i.name }

func (i fmvracer) Title() string {
	if i.qTime != "CHECK IN Please!" {
		name := i.name + " " + check
		return name
	} else {
		return i.name
	}

}
func (i fmvracer) Description() string { return i.qTime }
func (i fmvracer) FilterValue() string { return i.name }

func (m model) Init() tea.Cmd { return nil }

// Model Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		//filler until disable esc key
		case "tab":
			if m.state == fmvView {
				m.state = vdView

			} else {
				m.state = fmvView
				return m, cmd
			}
		case "enter":
			if m.state == vdView {
				index := m.fmv.Index()

				item := m.velocidrone.SelectedItem()
				cmd = m.fmv.InsertItem(index, item)
				cmds = append(cmds, cmd)
				m.fmv, cmd = m.fmv.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		case "r":
			if m.state == fmvView {
				index := m.fmv.Index()
				m.fmv.RemoveItem(index)
			}
		}
		switch m.state {
		case vdView:
			m.velocidrone, cmd = m.velocidrone.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case fmvView:
			m.fmv, cmd = m.fmv.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.velocidrone.SetSize(msg.Width-h, msg.Height-v)
		m.fmv.SetSize(msg.Width-h, msg.Height-v)
	}
	return m, tea.Batch(cmds...)

	//m.fmvRacers, cmd = m.fmvRacers.Update(msg)
	//m.racers, cmd = m.racers.Update(msg)
}

func (m model) View() string {
	m.velocidrone.DisableQuitKeybindings()
	m.fmv.DisableQuitKeybindings()

	left := docStyle.Render(m.velocidrone.View())
	right := docStyle.Render(m.fmv.View())
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	return body
}

func main() {
	vdList := []list.Item{}
	fmvList := []list.Item{}

	vdRacers := rt.GetVdRacers("race.csv")
	fmvRaw := rt.GetFMVvoice("checkin.csv")
	fmvRacers := rt.BindLists(vdRacers, fmvRaw)

	for _, v := range vdRacers {
		vdList = append(vdList, vdracer{name: v.VelocidronName, qTime: v.QualifyingTime, craft: v.ModelName})
	}

	for _, f := range fmvRacers {
		fmvList = append(fmvList, fmvracer{name: f.RacerName, qTime: f.QualifyingTime, craft: f.ModelName})
	}

	vItems := list.New(vdList, list.NewDefaultDelegate(), 0, 0)
	vItems.Styles.Title = itemStyle //wrong style name

	fmvItems := list.New(fmvList, list.NewDefaultDelegate(), 0, 0)
	fmvItems.Styles.Title = itemStyle //wrong style name

	m := model{velocidrone: vItems,
		fmv: fmvItems,
	}

	m.velocidrone.Title = "~Velocidrone Times~"
	m.fmv.Title = "~FMV Preflight Checkin~"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)

	}
}

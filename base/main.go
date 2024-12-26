package main

import (
	"fmt"
	"os"

	rt "abc.com/base/racetools"
	"github.com/charmbracelet/bubbles/key"
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
type listKeyMap struct {
	addRacer     key.Binding
	replaceRacer key.Binding
	removeRacer  key.Binding
	selectState  key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		addRacer: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add to FMV"),
		),

		replaceRacer: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r,", "replace FMV selection"),
		),

		removeRacer: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete racer"),
		),

		selectState: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch lists"),
		),
	}
}

const (
	vdView state = iota
	fmvView
)

type model struct {
	keys        *listKeyMap
	velocidrone list.Model
	fmv         list.Model
	state       state
}

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

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.velocidrone.SetSize(msg.Width-h, msg.Height-v)

		m.fmv.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.selectState):
			if m.state == fmvView {
				m.state = vdView
				m.velocidrone, cmd = m.velocidrone.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			if m.state == vdView {
				m.state = fmvView
				m.fmv, cmd = m.fmv.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)

			}
		}
		switch m.state {
		case vdView:
			if m.fmv.FilterState() == list.Filtering || m.velocidrone.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, m.keys.replaceRacer):
				index := m.fmv.Index()
				item := m.velocidrone.SelectedItem()
				Found := false
				for _, i := range m.fmv.Items() {
					if i == item {
						Found = true
					}
				}
				if !Found {
					m.fmv.RemoveItem(index)
					cmd = m.fmv.InsertItem(index, item)
					cmds = append(cmds, cmd)
				}
				m.fmv, cmd = m.fmv.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			case key.Matches(msg, m.keys.addRacer):
				item := m.velocidrone.SelectedItem()
				Found := false
				for _, i := range m.fmv.Items() {
					if i == item {
						Found = true
					}
				}
				if !Found {
					cmd = m.fmv.InsertItem(99999, item)
					cmds = append(cmds, cmd)
				}
				m.fmv, cmd = m.fmv.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		case fmvView:
			if m.fmv.FilterState() == list.Filtering || m.velocidrone.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, m.keys.removeRacer):
				index := m.fmv.Index()
				m.fmv.RemoveItem(index)
			}
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
			//filler until disable esc key
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

	}
	m.velocidrone, cmd = m.velocidrone.Update(msg)
	cmds = append(cmds, cmd)
	m.fmv, cmd = m.fmv.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	m.velocidrone.DisableQuitKeybindings()
	m.fmv.DisableQuitKeybindings()
	left := docStyle.Render(m.velocidrone.View())
	right := docStyle.Render(m.fmv.View())

	body := lipgloss.JoinHorizontal(lipgloss.Center, left, right)
	return body
}

func main() {
	listkeys := newListKeyMap()
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
	vItems.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listkeys.addRacer,
			listkeys.replaceRacer,
			listkeys.selectState,
		}
	}
	vItems.Styles.Title = itemStyle //wrong style name

	fmvItems := list.New(fmvList, list.NewDefaultDelegate(), 0, 0)
	fmvItems.Styles.Title = itemStyle //wrong style name
	fmvItems.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listkeys.removeRacer,
			listkeys.selectState,
		}
	}

	m := model{velocidrone: vItems,
		fmv:  fmvItems,
		keys: listkeys,
	}
	m.velocidrone.Title = "~Velocidrone Times~"
	m.fmv.Title = "~FMV Preflight Checkin~"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	rt "abc.com/base/racetools"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
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
	submitList   key.Binding
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

		submitList: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "Submit List"),
		),
	}
}

const (
	vdView state = iota
	fmvView
	tableView
)

type model struct {
	table           table.Model
	masterList      []*rt.Client
	checkedInRacers []string
	keys            *listKeyMap
	velocidrone     list.Model
	fmv             list.Model
	state           state
}

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	docStyle  = lipgloss.NewStyle().Margin(1, 2)
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("178")).
			Bold(true).
			Underline(true)
	magentaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("201")).
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
			case key.Matches(msg, m.keys.submitList):
				m.checkedInRacers = m.checkIn()
				list := m.addRacingList()
				rows := []table.Row{}
				for _, i := range list {
					rows = append(rows, i)
				}
				m.table.SetRows(rows)
				m.state = tableView
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
		case tableView:
			m.table, cmd = m.table.Update(msg)
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
	var body string
	switch {

	case m.state != tableView:
		m.velocidrone.DisableQuitKeybindings()
		m.fmv.DisableQuitKeybindings()
		left := docStyle.Render(m.velocidrone.View())
		right := docStyle.Render(m.fmv.View())

		body = lipgloss.JoinHorizontal(lipgloss.Center, left, right)
		return body
	case m.state == tableView:
		header := magentaStyle.Render("Magenta Group:")
		body = baseStyle.Render(m.table.View())
		body = lipgloss.JoinVertical(lipgloss.Center, header, body)
		return body
	}
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
			listkeys.submitList,
		}
	}
	columns := []table.Column{
		{Title: "Racer", Width: 10},
		{Title: "Time", Width: 6},
		{Title: "Quad", Width: 10},
	}
	rows := []table.Row{}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(8),
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

	m := model{velocidrone: vItems,
		fmv:        fmvItems,
		keys:       listkeys,
		masterList: vdRacers,
		table:      t,
	}
	m.velocidrone.Title = "~Velocidrone Times~"
	m.fmv.Title = "~FMV Preflight Checkin~"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) checkIn() []string {
	items := m.fmv.Items()
	var cList []string
	for _, i := range items {
		x := i.FilterValue()
		cList = append(cList, x)
	}
	return cList
}

func (m model) addRacingList() [][]string {
	type cleanRacer struct {
		racer string
		time  float64
		quad  string
	}

	var cleanRacers []cleanRacer

	var racers []*rt.Client

	for _, i := range m.checkedInRacers {
		for _, e := range m.masterList {
			if e.VelocidronName == i {
				racers = append(racers, e)
				break
			}
		}
	}

	for _, i := range racers {
		var x cleanRacer
		x.racer = i.VelocidronName
		x.quad = i.ModelName
		var t float64
		t, _ = strconv.ParseFloat(i.QualifyingTime, 64)
		x.time = t
		cleanRacers = append(cleanRacers, x)
	}

	sort.Slice(cleanRacers, func(i, j int) bool {
		return cleanRacers[i].time < cleanRacers[j].time
	})

	var racingList [][]string
	for _, i := range cleanRacers {
		var racestring []string
		racestring = append(racestring, i.racer)
		f := i.time
		s := strconv.FormatFloat(f, 'g', 5, 64)
		racestring = append(racestring, s)
		racestring = append(racestring, i.quad)
		racingList = append(racingList, racestring)
	}
	return racingList
}

package main

import (
	"fmt"
	"os"
	"sort"

	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	rt "abc.com/csvform/racetools"
)

const maxWidth = 80

var (
// red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
// indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
// green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("236"))
)

var fmvTag string = `



_____ __  ____     __                     
|  ___|  \/  \ \   / /                     
| |_  | |\/| |\ \ / /                      
|  _| | |  | | \ V /                       
|_|__ |_|  |_|  \_/  _   _ _ _       _ _ _ 
|   _ \ __ _ ___ ___| \ | (_) |_ ___| | | |
| |_)  / _ |/ __/ _ \  \| | | __/ _ \ | | |
|  _ < (_| | (_|  __/ |\  | | ||  __/_|_|_|
|_| \_\__,_|\___\___|_| \_|_|\__\___(_|_|_)

   
   `

type state uint

const (
	fmvState state = iota
	vdState
	formState
)

type Model struct {
	state   state
	csvForm *huh.Form
	width   int
	lg      *lipgloss.Renderer
	styles  *Styles

	discordList  []*rt.DiscordIds
	fmvVoiceList []*rt.FmvVoicePilot
	vdList       []*rt.VdPilot

	groups   []table.Model
	vdTable  table.Model
	fmvTable table.Model

	vdSearch list.Model
}

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)

	return &s
}

func (m Model) Init() tea.Cmd { return m.csvForm.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
		m.vdSearch.SetSize(20, 20)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "q":
			return m, tea.Quit
		case "`":
			m.csvForm.State = huh.StateNormal
			m.csvForm = NewModel().csvForm
			return m, cmd
		case "tab":
			if m.csvForm.State == huh.StateCompleted {
				if m.state == fmvState {
					m.state = vdState
				} else {
					m.state = fmvState
				}
			}
			return m, cmd

		case "A", "a":
			if m.state == fmvState {
				x := m.fmvTable.SelectedRow()
				m.Checkin(x)
				fmvRows := m.makeFMVTable()
				m.fmvTable.SetRows(fmvRows)
			}
			return m, cmd

		case "G", "g":
			var clearRows = []table.Row{}
			for i := 0; i < 5; i++ {
				m.groups[i].SetRows(clearRows)
			}
			list := m.addRacingList()      // order lists of entered
			brackets := rt.RaceArray(list) //
			//m.makeBrackets(brackets)       //allocates racers to groups
			indexLen := len(brackets)
			for i := 0; i < indexLen; i++ {
				rows := []table.Row{}
				for _, x := range brackets[i] {
					rows = append(rows, x)
					m.groups[i].SetRows(rows)
				}
			}
			for i := 1; i < 5; i++ {
				m.groups[i].Blur()

			}
		case "R", "r":
			m.VdUpdateFmvObj()
			fmvRows := m.makeFMVTable()
			m.fmvTable.SetRows(fmvRows)
		}

		switch m.state {
		case fmvState:
			m.fmvTable.Focus()
			m.fmvTable, cmd = m.fmvTable.Update(msg)
			cmds = append(cmds, cmd)
		case vdState:
			m.fmvTable.Blur()
		}
	}

	// Process the form
	form, cmd := m.csvForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.csvForm = f
		cmds = append(cmds, cmd)
	}
	if m.state == formState {
		if m.csvForm.State == huh.StateCompleted {

			m.discordList = rt.GetDiscordId(m.csvForm.GetString("discord"))
			m.fmvVoiceList = rt.GetFMVvoice(m.csvForm.GetString("fmv"))
			m.vdList = rt.GetVdRacers(m.csvForm.GetString("vd"))
			m.fmvVoiceList = rt.BindLists(m.vdList, m.fmvVoiceList, m.discordList)

			fmvrows := m.makeFMVTable()
			m.fmvTable.SetRows(fmvrows) //builds the fmv table from csv
			vdrows := m.makeVDTable()
			m.vdTable.SetRows(vdrows) // builds the vd table from vd

			m.vdSearch = m.makeList()

			m.fmvTable, cmd = m.fmvTable.Update(msg)
			m.state = vdState
			cmds = append(cmds, cmd)

		}
	}
	if m.state == vdState {
		m.vdSearch, cmd = m.vdSearch.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.csvForm.State {
	case huh.StateCompleted:

		header := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
		padding := lipgloss.NewStyle().Padding(0, 2)
		accii := lipgloss.NewStyle().Padding(0, 4).Foreground(lipgloss.Color("11"))
		tablePadding := lipgloss.NewStyle().Padding(1, 4)

		vdSearchBody := docStyle.Render(m.vdSearch.View())

		num := strconv.Itoa(len(m.fmvVoiceList))
		fmvtitle := header.Render(fmt.Sprintf("\n\n FMV Voice Checkin (count:%s)\n", num))
		fmvTable := m.fmvTable.View()
		fmvBody := padding.Render(lipgloss.JoinVertical(lipgloss.Center, fmvtitle, fmvTable))

		tables := lipgloss.JoinHorizontal(lipgloss.Center, vdSearchBody, fmvBody)
		fmvText := accii.Render(fmvTag)

		body := lipgloss.JoinHorizontal(lipgloss.Center, tables, fmvText)
		footer := "\nUse 'tab' to change lists\n"
		view := lipgloss.JoinVertical(lipgloss.Left, body, footer)

		//bracket groups section
		goldHeader := "Gold Group:"
		goldBody := baseStyle.Render(m.groups[0].View())
		gold := lipgloss.JoinVertical(lipgloss.Center, goldHeader, goldBody)

		mHeader := "Magenta Group:"
		mBody := baseStyle.Render(m.groups[1].View())
		magenta := lipgloss.JoinVertical(lipgloss.Center, mHeader, mBody)

		cyanHeader := "Cyan Group:"
		cyanBody := baseStyle.Render(m.groups[2].View())
		cyan := lipgloss.JoinVertical(lipgloss.Center, cyanHeader, cyanBody)

		orangeHeader := "Orange Group:"
		orangeBody := baseStyle.Render(m.groups[3].View())
		orange := lipgloss.JoinVertical(lipgloss.Center, orangeHeader, orangeBody)

		greenHeader := "Green Group:"
		greenBody := baseStyle.Render(m.groups[4].View())
		green := lipgloss.JoinVertical(lipgloss.Center, greenHeader, greenBody)

		r1s := lipgloss.JoinHorizontal(lipgloss.Center, gold, magenta, cyan, orange, green)
		r1 := tablePadding.Render(r1s)

		groupBody := lipgloss.JoinVertical(lipgloss.Center, view, r1)
		return groupBody

	default:
		body := m.csvForm.View()
		return body
	}

}

func main() {

	if _, err := tea.NewProgram(NewModel()).Run(); err != nil {
		fmt.Println("Bummer, there's been an error:", err)
		os.Exit(1)
	}
}

func NewModel() Model {

	m := Model{}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.csvForm = huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Velocidrone File").
				Description("CSV file for Veloidrone track").
				AllowedTypes([]string{".csv"}).
				Key("vd"),

			huh.NewFilePicker().
				Title("FMV Voice File").
				Description("CSV of FMV voice with User and ID flags").
				AllowedTypes([]string{".csv"}).
				Key("fmv"),

			huh.NewFilePicker().
				Title("Discord File").
				Description("CSV record of discord ID's and respective VD names").
				AllowedTypes([]string{".csv"}).
				Key("discord"),
		),
	).WithWidth(65).
		WithShowHelp(true).
		WithShowErrors(false)

	gColumns := []table.Column{
		{Title: "Name", Width: 11},
		{Title: "Time", Width: 6},
		{Title: "Craft", Width: 9},
	}
	vdColumns := []table.Column{
		{Title: "Name", Width: 16},
		{Title: "Time", Width: 8},
		{Title: "Craft", Width: 10},
	}

	fmvColumns := []table.Column{
		{Title: "Pilot", Width: 16},
		{Title: "VD Name", Width: 16},
		{Title: "Qualify time", Width: 16},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{}

	gt := table.New( //for color groups display
		table.WithColumns(gColumns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(8),
	)

	vdTable := table.New(
		table.WithColumns(vdColumns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(12),
	)

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
	vdTable.SetStyles(s)

	var groupTlist []table.Model
	group := table.DefaultStyles()
	group.Cell = group.Cell.
		UnsetForeground().
		Foreground(lipgloss.Color("172"))
	group.Selected = group.Selected.
		UnsetForeground().
		Foreground(lipgloss.Color("118"))
	gt.SetStyles(group)

	g1, g2, g3, g4, g5 := gt, gt, gt, gt, gt
	groupTlist = append(groupTlist, g1, g2, g3, g4, g5) // fix this ugly thing, just broken down to test

	var items = []list.Item{}
	m.vdSearch = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.vdSearch.Title = "Velocidrone Times"
	m.vdSearch.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("178")).
		Background(lipgloss.Color("0")).
		Bold(true).
		Underline(true)

	m.groups = groupTlist
	m.vdTable = vdTable
	m.fmvTable = fmvTable
	m.state = formState

	return m
}

func (m Model) makeFMVTable() []table.Row {
	rows := []table.Row{}

	for _, i := range m.fmvVoiceList {
		var s []string
		var fmvNul rt.FmvVoicePilot
		var status string
		name := i.RacerName
		vdName := i.VdName
		qtime := i.QualifyingTime
		if i.Status == fmvNul.Status {
			status = "-"
		} else {
			status = i.Status
		}
		s = append(s, name, vdName, qtime, status)
		rows = append(rows, s)
	}
	return rows
}
func (m Model) makeVDTable() []table.Row {
	rows := []table.Row{}
	for _, i := range m.vdList {
		var obj []string
		var name, time, quad string
		name = i.VelocidronName
		time = i.QualifyingTime
		quad = i.ModelName
		obj = append(obj, name, time, quad)
		rows = append(rows, obj)
	}
	return rows
}

var nulFmvVoice = rt.FmvVoicePilot{Status: "-"}

func (m Model) Checkin(r table.Row) {

	for _, i := range m.fmvVoiceList {
		if r[0] == i.RacerName {
			switch i.Status {
			case nulFmvVoice.VdName:
				if i.QualifyingTime != "CHECK IN Please!" {
					i.Status = "Entered"
				}
			case "Entered":
				i.Status = nulFmvVoice.VdName
			}
		}
	}
}

func (m Model) VdUpdateFmvObj() {
	r := m.fmvTable.SelectedRow()
	listItem := m.vdSearch.SelectedItem().FilterValue()
	for _, i := range m.fmvVoiceList {
		if r[0] == i.RacerName {
			for _, x := range m.vdList {
				if x.VelocidronName == listItem {
					i.VdName = x.VelocidronName
					i.QualifyingTime = x.QualifyingTime
					i.ModelName = x.ModelName
				}
			}
		}
	}
}

func (m Model) addRacingList() [][]string {
	type cleanRacer struct {
		racer string
		time  float64
		quad  string
	}

	var cleanRacers []cleanRacer
	var racers []*rt.FmvVoicePilot

	for _, i := range m.fmvVoiceList {
		if i.Status == "Entered" {
			racers = append(racers, i)
		}
	}

	for _, i := range racers {
		var x cleanRacer
		x.racer = i.VdName
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

type item struct {
	name, time, craft string
}

func (i item) Title() string { return i.name }
func (i item) Description() string {
	description := fmt.Sprintf("%s | %s", i.time, i.craft)
	return description
}
func (i item) FilterValue() string { return i.name }

var docStyle = lipgloss.NewStyle().Padding(7, 14, 0, 0) // for list

func (m Model) makeList() list.Model {
	for _, i := range m.vdList {
		obj := item{name: i.VelocidronName, time: i.QualifyingTime, craft: i.ModelName}
		//items = append(items, obj)
		m.vdSearch.InsertItem(99999, obj)
	}
	return m.vdSearch
}

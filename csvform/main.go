package main

import (
	"fmt"
	"os"

	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	rt "abc.com/csvform/racetools"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
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
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

func (m Model) Init() tea.Cmd { return m.csvForm.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()

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

		}

		switch m.state {
		case fmvState:
			m.vdTable.Blur()
			m.fmvTable.Focus()
			m.fmvTable, cmd = m.fmvTable.Update(msg)
			cmds = append(cmds, cmd)

		case vdState:
			m.vdTable.Focus()
			m.vdTable, cmd = m.vdTable.Update(msg)
			m.fmvTable.Blur()
			cmds = append(cmds, cmd)

			//return m, tea.Batch(cmds...)
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

			m.fmvTable, cmd = m.fmvTable.Update(msg)
			m.state = fmvState
			cmds = append(cmds, cmd)

		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.csvForm.State {
	case huh.StateCompleted:
		header := lipgloss.NewStyle().Foreground(lipgloss.Color("44"))
		padding := lipgloss.NewStyle().Padding(0, 2)
		accii := lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("184"))

		vdTitle := header.Render("\n\nVelocidrone Time\n")
		vdTable := m.vdTable.View()
		vdBody := padding.Render(lipgloss.JoinVertical(lipgloss.Center, vdTitle, vdTable))

		num := strconv.Itoa(len(m.fmvVoiceList))
		fmvtitle := header.Render(fmt.Sprintf("\n\n FMV Voice Checkin (count:%s)\n", num))
		fmvTable := m.fmvTable.View()
		fmvBody := padding.Render(lipgloss.JoinVertical(lipgloss.Center, fmvtitle, fmvTable))

		tables := lipgloss.JoinHorizontal(lipgloss.Top, vdBody, fmvBody)
		fmvText := accii.Render(fmvTag)

		body := lipgloss.JoinHorizontal(lipgloss.Center, tables, fmvText)
		footer := "\n\nUse 'tab' to change lists"
		view := lipgloss.JoinVertical(lipgloss.Left, body, footer)
		return view

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
		{Title: "Name", Width: 16},
		{Title: "Time", Width: 8},
	}
	vdColumns := []table.Column{
		{Title: "Name", Width: 16},
		{Title: "Time", Width: 8},
		{Title: "Craft", Width: 10},
	}

	fmvColumns := []table.Column{
		{Title: "Name", Width: 16},
		{Title: "Qualify time", Width: 14},
		{Title: "Status", Width: 10},
	}

	rows := []table.Row{}

	groupTable := table.New( //for color groups display
		table.WithColumns(gColumns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(6),
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
		BorderForeground(lipgloss.Color("3")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	var groupTlist []table.Model

	for i := 0; i < 5; i++ {
		groupTlist = append(groupTlist, groupTable)
	}
	fmvTable.SetStyles(s)
	vdTable.SetStyles(s)

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
		qtime := i.QualifyingTime
		if i.Status == fmvNul.Status {
			status = "Missing"
		}
		s = append(s, name, qtime, status)
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

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

type Model struct {
	csvForm *huh.Form
	width   int
	lg      *lipgloss.Renderer
	styles  *Styles
}

func NewModel() Model {
	m := Model{}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	var (
		VdFile       string // csv of the weekly track and qualifying times.
		FmvVoiceFile string // voice chat csv with usernames and discord IDs.
		DiscordFile  string // file that has racer's IDs paired to velocidrone names.
	)

	m.csvForm = huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Velocidrone File").
				Description("CSV file for Veloidrone track").
				AllowedTypes([]string{".csv"}).
				Value(&VdFile),

			huh.NewFilePicker().
				Title("FMV Voice File").
				Description("CSV of FMV voice with User and ID flags").
				AllowedTypes([]string{".csv"}).
				Value(&FmvVoiceFile),

			huh.NewFilePicker().
				Title("Discord File").
				Description("CSV record of discord ID's and respective VD names").
				AllowedTypes([]string{".csv"}).
				Value(&DiscordFile),
		),
	).WithWidth(65).
		WithShowHelp(true).
		WithShowErrors(false)
	return m
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
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.csvForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.csvForm = f
		cmds = append(cmds, cmd)
	}

	if m.csvForm.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.csvForm.State {
	case huh.StateCompleted:
		body := "OOOO babyyyyy"
		return body
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

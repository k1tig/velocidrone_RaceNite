package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/gocarina/gocsv"
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

//type state int

const (
// statusNormal state = iota
// stateDone
)

type Model struct {
	//state  state
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
	//VD filtered list of times
	racers list.Model
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("fmvname").
				Title("FMV Name").
				Description("Enter FMV Discord Name"), // test names for validation

			huh.NewInput().
				Key("vdname").
				Title("Velocidrone Name").
				Description("Enter 'SolaFide' or 'jon E5'"),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					if m.form.GetString("fmvname") == "" {
						return fmt.Errorf("Enter Missing Fields")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m Model) Init() tea.Cmd {
	m.getVDsheet()
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	switch m.form.State {

	//Racer entry confirmation
	case huh.StateCompleted:

		var b strings.Builder
		//Send VD and Discord names to raceNite Master list * not made yet
		//toCheckinList := m.form.GetString("fmvname")
		fmt.Fprintf(&b, "Racer Succefully Entered: %s", m.form.GetString("fmvname"))
		return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
	default:

		var fmvname string

		if m.form.GetString("fmvname") != "" {
			fmvname = "FMV Name: " + m.form.GetString("fmvname")
		}

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			var (
				buildInfo = "Waiting..."
				vdinfo    string
			)

			if m.form.GetString("fmvname") != "" && m.form.GetString("vdname") != "" {
				vdinfo = m.getVdUser()
				vdinfo = "\n\n" + s.StatusHeader.Render("Velocidrone User Info") + "\n" + vdinfo
				buildInfo = fmt.Sprintf("%s\n", fmvname)

			}

			const statusWidth = 28
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("FMV Discord Checkin") + "\n" +
					buildInfo +
					vdinfo)
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("FMV RaceNite Manual Entry Form, Baybeeeeeeee!!!")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m Model) getVdUser() string {
	for _, racer := range okRaceClass {
		formName := m.form.GetString("vdname")
		if racer.PlayerName == formName {
			vdUserMatch := ("FOUND...\n\nUsername: " + racer.PlayerName + "\nLaptime: " + racer.LapTime + "\nCraft: " + racer.ModelName)
			return vdUserMatch
		}
	}
	return "Racer not found"
}

func main() {

	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}

// Generic Data

type Client struct { // Our example struct, you can use "-" to ignore a field
	PlayerName string `csv:"Player Name"`
	LapTime    string `csv:"Lap Time"`
	X_Pos      string `csv:"-"`
	ModelName  string `csv:"Model Name"`
	X_Country  string `csv:"-"`
}

var clients = []*Client{}
var okRaceClass = []*Client{}
var vdList = []list.Item{}

type Vdracer struct {
	name, qualTime, craft string
}

func (i Vdracer) Title() string       { return i.name }
func (i Vdracer) Description() string { return i.qualTime + " | " + i.craft }
func (i Vdracer) FilterValue() string { return i.name }

func (m Model) getVDsheet() tea.Model {
	raceFile, err := os.OpenFile("race.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer raceFile.Close()
	if err := gocsv.UnmarshalFile(raceFile, &clients); err != nil { // Load clients from file
		panic(err)
	}
	for _, client := range clients { //clients are the master qual times
		if client.ModelName == "TBS Spec" || client.ModelName == "Twig XL 3" {
			okRaceClass = append(okRaceClass, client) // checkedIn seperates the class of quads from the master list
		}
	}
	// Imports the VD csv to the vdList
	for _, r := range okRaceClass {
		vdList = append(vdList, Vdracer{name: r.PlayerName, qualTime: r.LapTime, craft: r.ModelName})
	}
	// Adds the list to the model. Important structure to allow for future styling
	r := list.New(vdList, list.NewDefaultDelegate(), 0, 0)
	m.racers = r

	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return m
}

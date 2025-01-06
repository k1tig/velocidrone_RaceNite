package main

import (
	"fmt"
	"os"

	rt "abc.com/sandbox/racetools"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

const maxWidth = 80


// use websocket as msg system to invoke func for checking Google sheets master copy
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

type Model struct {
	Split_1 table.Model
	racers  []racer
	form    *huh.Form

	width  int
	lg     *lipgloss.Renderer
	styles *Styles
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

func initModel() Model {

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

	m := Model{width: maxWidth}

	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.racers = rList
	m.Split_1 = t
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("class").
				Options(huh.NewOptions("Gold", "Magenta", "Teel")...).
				Title("Race Group").
				Description("Race bracket to be entered"),

			huh.NewSelect[string]().
				Key("level").
				Options(huh.NewOptions("asiy: 60718", "eedok: 69", "uGellin: 61.718", "anddy: 64.551")...).
				Title("Heats").
				Description("How slow you were last round"),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
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
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
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
		case "e":
			vdRacers := rt.GetVdRacers("race.csv")
			for e, i := range m.racers {
				for _, x := range vdRacers {
					if i.name == x.VelocidronName {
						var nulCheck racer // represents empty value of type racer.B1 (qualifying time 1)
						if m.racers[e].B1 == nulCheck.B1 {
							if x.ModelName == "TBS Spec" || x.ModelName == "Twig XL 3" {
								m.racers[e].B1 = x.QualifyingTime
								break
							}

						}
					}
				}
			}

			x := m.buildTable()
			m.Split_1.SetRows(x)
		}

	}
	m.Split_1, cmd = m.Split_1.Update(msg)
	return m, cmd
}
func (m Model) View() string {
	v := m.form.View()
	form := m.lg.NewStyle().Margin(0, 4).Render(v)
	table := baseStyle.Render(m.Split_1.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, table, form)
	return body
}

func main() {
	if _, err := tea.NewProgram(initModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m Model) buildTable() []table.Row {
	newRows := []table.Row{}
	for _, i := range m.racers {
		x := []string{i.name, i.B1, i.r1, i.p1, i.r2, i.p2, i.r3, i.p3}
		newRows = append(newRows, x)

	}

	return newRows

}

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

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

//Function to update racers from a list

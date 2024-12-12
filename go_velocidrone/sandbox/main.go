package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gocarina/gocsv"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// represents VD racers qualifying times
type vdracer struct {
	name, qualTime, craft string
}

// represents FMV racers checked in raceNite
type fmvracer struct {
	name, qualTime, craft string
}

func (i vdracer) Title() string       { return i.name }
func (i vdracer) Description() string { return i.qualTime + " | " + i.craft }
func (i vdracer) FilterValue() string { return i.name }

type model struct {
	racers list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.racers.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.racers, cmd = m.racers.Update(msg)
	return m, cmd
}

func (m model) View() string {
	left := docStyle.Render(m.racers.View())
	right := docStyle.Render(m.racers.View())
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	return body
}

func main() {
	getRacers()
	items := []list.Item{}
	for _, r := range checkedIn {
		items = append(items, vdracer{name: r.PlayerName, qualTime: r.LapTime, craft: r.ModelName})
	}

	m := model{racers: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.racers.Title = "Add Racers to Sheets:"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// compontents for retrieving VD qualifying time sheet
type Client struct { // Our example struct, you can use "-" to ignore a field
	PlayerName string `csv:"Player Name"`
	LapTime    string `csv:"Lap Time"`
	X_Pos      string `csv:"-"`
	ModelName  string `csv:"Model Name"`
	X_Country  string `csv:"-"`
}

var clients = []*Client{}
var checkedIn = []*Client{}

type racer struct {
	name string
}

func getRacers() {
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
			checkedIn = append(checkedIn, client) // checkedIn seperates the class of quads from the master list
		}
	}

	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
}

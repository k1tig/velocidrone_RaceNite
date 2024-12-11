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

type item struct {
	name, raceTime, craft string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.raceTime + " | " + i.craft }
func (i item) FilterValue() string { return i.name }

type model struct {
	list list.Model
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
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	getRacers()
	items := []list.Item{}
	for _, r := range checkedIn {
		items = append(items, item{name: r.PlayerName, raceTime: r.LapTime, craft: r.ModelName})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Add Racers to Sheets:"

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

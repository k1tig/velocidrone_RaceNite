package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gocarina/gocsv"
)

var (
	docStyle  = lipgloss.NewStyle().Margin(1, 2)
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("178")).
			Background(lipgloss.Color("0")).
			Bold(true).
			Underline(true)
)

// represents VD racers qualifying times
type vdracer struct {
	name, qualTime, craft string
}

// represents FMV racers checked in raceNite
type fmvracer struct {
	name string
}

func (i vdracer) Title() string       { return i.name }
func (i vdracer) Description() string { return i.qualTime + " | " + i.craft }
func (i vdracer) FilterValue() string { return i.name }

func (i fmvracer) Title() string       { return i.name }
func (i fmvracer) Description() string { return "" }
func (i fmvracer) FilterValue() string { return i.name }

type state uint

const (
	vdView state = iota
	fmvView
)

// Model Struct
type model struct {
	racers    list.Model
	fmvRacers list.Model
	state     state
}

// Moedl Init
func (m model) Init() tea.Cmd {
	return nil
}

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
		}
		switch m.state {
		case vdView:
			m.racers, cmd = m.racers.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case fmvView:
			m.fmvRacers, cmd = m.fmvRacers.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.racers.SetSize(msg.Width-h, msg.Height-v)
		m.fmvRacers.SetSize(msg.Width-h, msg.Height-v)
	}
	return m, tea.Batch(cmds...)

	//m.fmvRacers, cmd = m.fmvRacers.Update(msg)
	//m.racers, cmd = m.racers.Update(msg)
}

// Model View
func (m model) View() string {
	m.racers.DisableQuitKeybindings()
	m.fmvRacers.DisableQuitKeybindings()

	left := docStyle.Render(m.racers.View())
	right := docStyle.Render(m.fmvRacers.View())
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	return body
}

// Main
func main() {
	vdList := []list.Item{}
	checkInList := []list.Item{}
	getRacers()
	getFMVvoice()
	voiceCheckinRacers()

	for _, r := range okRaceClass {
		vdList = append(vdList, vdracer{name: r.PlayerName, qualTime: r.LapTime, craft: r.ModelName})
	}

	for _, checkedR := range checkedInRacers {
		checkInList = append(checkInList, fmvracer{name: checkedR.Racer})
	}

	getFMVvoice()

	//need to restructure making lists to aaccount for styling
	r := list.New(vdList, list.NewDefaultDelegate(), 0, 0)
	r.Styles.Title = itemStyle //wrong style name

	f := list.New(checkInList, list.NewDefaultDelegate(), 0, 0)
	f.Styles.Title = itemStyle //wrong style name

	m := model{racers: r,
		fmvRacers: f,
	}

	m.racers.Title = "~Velocidrone Times~"
	m.fmvRacers.Title = "~Checked in~"

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

// Yoinks FMV Voice chat members
type Racers struct {
	Racer string `csv:"Display Name"`
}

var clients = []*Client{}
var okRaceClass = []*Client{}

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
			okRaceClass = append(okRaceClass, client) // checkedIn seperates the class of quads from the master list
		}
	}
	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
}

var fmvRacers = []*Racers{}

func getFMVvoice() {
	fmvVoiceFile, err := os.OpenFile("checkin.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fmvVoiceFile.Close()
	if err := gocsv.UnmarshalFile(fmvVoiceFile, &fmvRacers); err != nil { // Load clients from file
		panic(err)
	}
	if _, err := fmvVoiceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
}

var checkedInRacers = []*Racers{}

// adds racers from fmv voice chat with a qualifying time + class quad OK to Checked in List
func voiceCheckinRacers() {
	for _, fmvR := range fmvRacers {
		for _, vdR := range okRaceClass {
			if vdR.PlayerName == fmvR.Racer {
				checkedInRacers = append(checkedInRacers, fmvR)
				break
			}
		}
	}
}

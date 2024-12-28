package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gocarina/gocsv"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table []table.Model
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		rows := addRows()
		m.table[1].SetRows(rows[0])
		m.table[0].SetRows(rows[1])
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table[1].Focused() {
				m.table[1].Blur()
			} else {
				m.table[1].Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table[1].SelectedRow()[1]),
			)
		}
	}
	m.table[1], cmd = m.table[1].Update(msg)
	return m, cmd
}

func (m model) View() string {
	leftText := "Magenta Group"
	left := baseStyle.Render(m.table[1].View())
	left = lipgloss.JoinVertical(lipgloss.Center, leftText, left)
	rightText := "Gold Group"
	right := baseStyle.Render(m.table[0].View())
	right = lipgloss.JoinVertical(lipgloss.Center, rightText, right)
	body := lipgloss.JoinHorizontal(lipgloss.Left, left, right)
	return body
}

func main() {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}
	rows := []table.Row{}
	t1 := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	t2 := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
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
	t1.SetStyles(s)
	t2.SetStyles(s)
	type tables []table.Model
	m := model{tables{t1, t2}}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func addRows() [][]table.Row {
	var listOfRows [][]table.Row
	rows1 := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
		{"7", "Cairo", "Egypt", "21,750,020"},
	}
	rows2 := []table.Row{
		{"8", "a", "Japan", "37,274,000"},
		{"9", "b", "India", "32,065,760"},
		{"10", "c", "China", "28,516,904"},
		{"11", "D", "Bangladesh", "22,478,116"},
		{"12", "ã Paulo", "Brazil", "22,429,800"},
		{"13", "M", "Mexico", "22,085,140"},
		{"14", "z", "Egypt", "21,750,020"},
	}
	listOfRows = append(listOfRows, rows1, rows2)
	return listOfRows
}

type Client struct { //struct to recieve data from velocidrone csv
	VelocidronName string `csv:"Player Name"`
	QualifyingTime string `csv:"Lap Time"`
	X_Pos          string `csv:"-"`
	ModelName      string `csv:"Model Name"`
	X_Country      string `csv:"-"`
}

type Racers struct {
	RacerName      string `csv:"Display Name"`
	VelocidronName string
	QualifyingTime string
	ModelName      string
}

// take a list of racers and returns group sets of racers

func GetVdRacers(filename string) []*Client {

	var Clients = []*Client{}
	var OkRaceClass = []*Client{}

	raceFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer raceFile.Close()
	if err := gocsv.UnmarshalFile(raceFile, &Clients); err != nil { // Load clients from file
		panic(err)
	}
	for _, client := range Clients { //clients are the master qual times
		if client.ModelName == "TBS Spec" || client.ModelName == "Twig XL 3" {
			OkRaceClass = append(OkRaceClass, client) // checkedIn seperates the class of quads from the master list
		}
	}
	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return OkRaceClass
}

// For using the voice chat in FMV discord as base group for pairing.
func GetFMVvoice(fileCsv string) []*Racers {

	var FmvRacers = []*Racers{}

	fmvVoiceFile, err := os.OpenFile(fileCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fmvVoiceFile.Close()
	if err := gocsv.UnmarshalFile(fmvVoiceFile, &FmvRacers); err != nil { // Load clients from file
		fmt.Printf("Something broke with FMV CSV: %v", err) //csv needs to be in same folder as main.go for now
	}
	if _, err := fmvVoiceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return FmvRacers
}

func RaceArray(vdList []*Client) [][]*Client {
	var maxGroupsize = 8
	var grouplength int
	var totalGroups int
	var modulus int

	racers := len(vdList)
	if racers > 40 {
		maxGroupsize = 10
	}

	for i := 1; i <= maxGroupsize; i++ {
		if racers/i <= maxGroupsize {
			totalGroups = i
			modulus = racers % i
			if modulus == 0 {
				grouplength = racers / i
			} else {
				grouplength = (racers - modulus) / i
			}
			break
		}
	}

	var groupStructure = make([][]*Client, totalGroups)
	var c int
	x := modulus

	for i := 1; i <= totalGroups; i++ {
		if x > 0 { // distribues the modulus between the lower teir groups
			racers := vdList[c : i*(grouplength+1)]
			groupStructure[i-1] = racers
			x--
			c += grouplength + 1
		} else { // groups that don't take a modulus
			racers := vdList[c : c+grouplength]
			groupStructure[i-1] = racers
			c += grouplength
		}
	}
	return groupStructure
}

func BindLists(vdl []*Client, fmvl []*Racers) []*Racers {
	var bound []*Racers
	for _, f := range fmvl {
		for _, v := range vdl {
			if v.VelocidronName == f.RacerName {
				f.VelocidronName = v.VelocidronName
				f.QualifyingTime = v.QualifyingTime
				f.ModelName = v.ModelName
				bound = append(bound, f)
				break
			}

		}
		if f.VelocidronName == "" {
			f.QualifyingTime = "CHECK IN Please!"
			bound = append(bound, f)
		}
	}
	return bound
}

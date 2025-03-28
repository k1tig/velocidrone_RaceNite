package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const maxWidth = 80

var rooms = []string{"1", "3", "5", "7"}

var currentBracket string

type Model struct {
	width int
	form  *huh.Form
}

func NewModel() Model {
	m := Model{width: maxWidth}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("raceid").
				Title("Active Races:").
				OptionsFunc(func() []huh.Option[string] {
					s := rooms
					// simulate API call
					time.Sleep(1 * time.Second)
					return huh.NewOptions(s...)
				},
					huh.NewConfirm().
						Key("done").
						Title("All done?").
						Validate(func(v bool) error {
							if !v {
								return fmt.Errorf("finish it")
							}
							return nil
						}).
						Affirmative("Yep").
						Negative("Wait, no"),
				),
		).WithWidth(45).
			WithShowHelp(false).
			WithShowErrors(false),
	)
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
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
	}

	return m, tea.Batch(cmds...)
}
func (m Model) View() string {
	switch m.form.State {
	case huh.StateCompleted:
		currentBracket = m.form.GetString("raceid")
		return "Room id: " + currentBracket
	default:
		form := m.form.View()
		return form
	}
}

func main() {
	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}

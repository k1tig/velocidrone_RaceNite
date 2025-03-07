package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var tui *Tui

func main() {
	tui = NewTui()
	p := tea.NewProgram(tui, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

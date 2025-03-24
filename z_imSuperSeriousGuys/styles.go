package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func fmvTableSelectedStyle(bgColor, fgColor string) table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("ffb3fd")).
		Foreground(lipgloss.Color("239")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor))

	return s
}

func vdSearchSelectedStyle(color string) list.DefaultDelegate {
	s := list.NewDefaultDelegate()
	s.Styles.SelectedTitle = s.Styles.SelectedTitle.
		Foreground(lipgloss.Color(color)).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: color})
	s.Styles.SelectedDesc = s.Styles.SelectedTitle

	return s
}

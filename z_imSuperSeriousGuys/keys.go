package main

import "github.com/charmbracelet/bubbles/key"

type fmvTableKeyMap struct {
	checkin    key.Binding
	checkinAll key.Binding
	remove     key.Binding
}

func (k fmvTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.checkin, k.checkinAll, k.remove}
}

func (k fmvTableKeyMap) FullHelp() []key.Binding {
	return []key.Binding{k.checkin, k.checkinAll, k.remove}
}

var theFmvKeys = fmvTableKeyMap{
	checkin: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Checkin"),
	),
	checkinAll: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "Checkin-all"),
	),
	remove: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "Delete Pilot"),
	),
}

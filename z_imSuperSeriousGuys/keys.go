package main

import "github.com/charmbracelet/bubbles/key"

type fmvTableKeyMap struct {
	checkin    key.Binding
	checkinAll key.Binding
	remove     key.Binding
	switchToVd key.Binding
}

func (k fmvTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.switchToVd, k.checkinAll, k.checkin, k.remove}
}

func (k fmvTableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.switchToVd, k.checkinAll, k.checkin, k.remove},
	}
}

var theFmvKeys = fmvTableKeyMap{
	checkinAll: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "Checkin-all"),
	),
	checkin: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Checkin"),
	),
	remove: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "Delete Pilot"),
	),
	switchToVd: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Vecocidrone List"),
	),
}

type vdSearchKeyMap struct {
	addToFmv    key.Binding
	updateAtFmv key.Binding
	switchToFmV key.Binding
}

func (k vdSearchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.switchToFmV, k.addToFmv, k.updateAtFmv}
}

func (k vdSearchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.switchToFmV, k.addToFmv, k.updateAtFmv}}
}

var theVdSearchKeys = vdSearchKeyMap{
	addToFmv: key.NewBinding(
		key.WithKeys("A", "a"),
		key.WithHelp("A/a", "Add to FMV list"),
	),
	updateAtFmv: key.NewBinding(
		key.WithKeys("U"),
		key.WithHelp("U", "Update FMV-Pilot"),
	),
	switchToFmV: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Swith to FMV table"),
	),
}

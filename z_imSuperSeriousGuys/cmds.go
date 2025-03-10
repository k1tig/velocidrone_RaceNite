package main

import tea "github.com/charmbracelet/bubbletea"

type csvProcessedMsg [][]Pilot

func (e csvForm) csvProcessedCmd() tea.Cmd {
	return func() tea.Msg {
		return e.processForm()
	}
}
func (e csvForm) processForm() csvProcessedMsg {
	discordTarget := GetDiscordId(e.form.GetString("discord"))
	fmvTarget := GetFMVvoice(e.form.GetString("fmv"))
	vdTarget := GetVdRacers(e.form.GetString("vd"))
	registeredTarget := BindLists(vdTarget, fmvTarget, discordTarget)
	return csvProcessedMsg{discordTarget, fmvTarget, vdTarget, registeredTarget}

}

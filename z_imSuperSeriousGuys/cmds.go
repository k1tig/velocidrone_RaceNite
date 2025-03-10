package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type csvProcessedMsg [][]Pilot

func processForm(e *huh.Form) (d1, f1, v1, reg1 []Pilot) {
	discordTarget := GetDiscordId(e.GetString("discord"))
	fmvTarget := GetFMVvoice(e.GetString("fmv"))
	vdTarget := GetVdRacers(e.GetString("vd"))
	registeredTarget := BindLists(vdTarget, fmvTarget, discordTarget)
	return discordTarget, fmvTarget, vdTarget, registeredTarget

}

type testMsg struct{}

func testCmd() tea.Msg {
	return testMsg{}
}

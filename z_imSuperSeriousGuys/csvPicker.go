package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type formMsg struct{}

func formCmd() tea.Cmd {
	return func() tea.Msg {
		return formMsg{}
	}
}

type mainMsg struct{}

func mainViewMsg() tea.Msg {
	return mainMsg{}
}

type csvForm struct {
	form      *huh.Form
	formReady bool
}

func (e csvForm) Init() tea.Cmd { return e.form.Init() }
func (e csvForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return e, nil
	case formMsg:
		e = initForm()
		e.form.Init()
		e.formReady = true

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return e, tea.Quit
		}
	}
	if e.formReady { //to keep from trying to update the form before init()?
		form, cmd := e.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			e.form = f
			cmds = append(cmds, cmd)
		}

		if e.form.State == huh.StateCompleted {
			return tui.Update(mainViewMsg())
		}

	}

	cmds = append(cmds, cmd)
	return e, tea.Batch(cmds...)
}

func (e csvForm) View() string {
	if e.formReady {
		return e.form.View()
	}
	return "Form Not Generated"
}

func initForm() csvForm {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Velocidrone File").
				Description("CSV file for Veloidrone track").
				AllowedTypes([]string{".csv"}).
				Key("vd"),

			huh.NewFilePicker().
				Title("FMV Voice File").
				Description("CSV of FMV voice with User and ID flags").
				AllowedTypes([]string{".csv"}).
				Key("fmv"),

			huh.NewFilePicker().
				Title("Discord File").
				Description("CSV record of discord ID's and respective VD names").
				AllowedTypes([]string{".csv"}).
				Key("discord"),
		),
	).WithWidth(65).
		WithShowHelp(true).
		WithShowErrors(false)
	f := csvForm{form: form}
	return f
}

package main

//This will be the main routing file for all internal operations
//such as those which don't create external services

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	mainView viewState = iota
	csvView
	createView
	observeView
	modView
	settingsView
	testView
)

type Tui struct {
	createForm entryForm
	list       list.Model
	state      viewState
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type mi string                 // front page menu link
func (mi) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(mi)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

func NewTui() *Tui {
	const defaultWidth = 20
	menuItems := []list.Item{
		mi("Create Race"),
		mi("Spectate"),
		mi("Moderate Race"),
		mi("Help / Settings"),
	}
	l := list.New(menuItems, itemDelegate{}, defaultWidth, 14)
	l.Title = "FMV RaceNite"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &Tui{list: l, state: mainView}
}

func (m Tui) Init() tea.Cmd {
	return nil
}

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch m.state {
		case mainView:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				switch m.list.SelectedItem().(mi) {
				case "Create Race":
					m.state = createView
					m, cmd := m.createForm.Update(msg)
					cmds = append(cmds, cmd, formCmd())
					return m, tea.Batch(cmds...)
				}
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	case mainMsg:
		m.state = testView
	case entryForm:
	}
	return m, tea.Batch(cmds...)
}

func (m Tui) View() string {
	if m.state != mainView {
		switch m.state {
		case csvView:
		case testView:
			return "\n" + "Booooooty"
		}
	}
	return "\n" + m.list.View()
}

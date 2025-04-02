package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

const gap = "\n\n"

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error

	done chan struct{}
	msg  chan []byte
	conn *websocket.Conn
}
type initReaderMsg struct{}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	msg := make(chan []byte)
	done := make(chan struct{})

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		msg:         msg,
		done:        done,
		conn:        c,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.initReader(), waitForMsg(m.msg))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			err := m.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Printf("error: %s", err)
			}

			return m, tea.Quit
		case tea.KeyEnter:
			text := m.textarea.Value()
			err := m.conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Fatal(err)
			}
		}
	case responseMsg:
		m.messages = append(m.messages, m.senderStyle.Render("You: ")+string(msg))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForMsg(m.msg), m.initReader())

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	cmds = append(cmds, tiCmd, vpCmd)
	return m, tea.Batch(cmds...)
}
func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
func (m model) initReader() tea.Cmd {
	return func() tea.Msg {
		defer close(m.done)
		for {
			_, message, err := m.conn.ReadMessage()
			if err != nil {
				fmt.Printf("err: %s", err)
			}
			m.msg <- message
		}
	}
}

type responseMsg []byte

func waitForMsg(sub chan []byte) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

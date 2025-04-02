package tui

import (
	"log/slog"
	"p1/pkg/api"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Splash struct {
	data  bool
	delay bool
}

type DelayCompleteMsg struct{}

type BrokerDataLoaded struct{}

type WebSocketConnected struct{}

func (m model) LoadCmds() []tea.Cmd {
	cmds := []tea.Cmd{}

	// Make sure the loading state shows for at least a couple seconds
	cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return DelayCompleteMsg{}
	}))

	// Initialize websocket client
	m.wsClient = NewWebSocketClient("ws://localhost:8080/ws") // Replace with your actual websocket URL
	cmds = append(cmds, m.wsClient.Connect())

	return cmds
}

func (m model) IsLoadingComplete() bool {
	return m.splash.data &&
		m.splash.delay &&
		m.wsConnected
}

func (m model) SplashInit() tea.Cmd {
	cmd := func() tea.Msg {
		api.FetchBrokerState()
		return BrokerDataLoaded{}
	}
	disableMouseCmd := func() tea.Msg {
		return tea.DisableMouse()
	}

	return tea.Batch(m.CursorInit(), disableMouseCmd, cmd)
}

func (m model) SplashUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case BrokerDataLoaded:
		m.splash.data = true
		return m, tea.Batch(m.LoadCmds()...)
	case DelayCompleteMsg:
		m.splash.delay = true
	case WebSocketUpdateMsg:
		if data, ok := msg.Data.(map[string]interface{}); ok {
			if data["type"] == "connected" {
				m.wsConnected = true
				return m, nil
			}
			// Handle other message types here
			slog.Info("Received websocket update", "data", data)
		}
	}

	if m.IsLoadingComplete() {
		return m.InitialDataLoaded()
	}
	return m, nil
}

func (m model) SplashView() string {
	var msg string
	if m.error != nil {
		msg = m.error.message
	} else {
		msg = ""
	}

	var hint string
	if m.error != nil {
		hint = lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.theme.TextAccent().Bold(true).Render("esc"),
			" ",
			"quit",
		)
	} else {
		hint = ""
	}

	if m.error == nil {
		return lipgloss.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			m.LogoView(),
		)
	}

	return lipgloss.Place(
		m.viewportWidth,
		m.viewportHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			"",
			"",
			"",
			m.LogoView(),
			"",
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.theme.TextError().Render(msg),
			),
			hint,
		),
	)
}

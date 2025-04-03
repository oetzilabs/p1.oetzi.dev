package models

import (
	"log/slog"
	"p1/pkg/api"
	"p1/pkg/tui/theme"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Splash struct {
	data     bool
	delay    bool
	theme    *theme.Theme
	wsClient *api.WebSocketClient
	logo     *Logo
	width    int
	height   int
}

type DelayCompleteMsg struct{}

type BrokerDataLoaded struct{}

type WebSocketConnected struct{}

func (sp *Splash) LoadCmds() []tea.Cmd {
	cmds := []tea.Cmd{}
	// Make sure the loading state shows for at least a couple seconds
	cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return DelayCompleteMsg{}
	}))

	// cmds = append(cmds, sp.wsClient.Connect())

	return cmds
}

func NewSplash(theme *theme.Theme, wsClient *api.WebSocketClient) *Splash {
	cursor := NewCursor(theme)
	return &Splash{
		theme:    theme,
		wsClient: wsClient,
		logo:     NewLogo(theme, cursor),
	}
}

func (sp *Splash) Init() tea.Cmd {
	slog.Info("Initializing Splash")

	cmd := func() tea.Msg {
		time.Sleep(time.Second * 1)
		return BrokerDataLoaded{}
	}

	disableMouseCmd := func() tea.Msg {
		return tea.DisableMouse()
	}
	return tea.Batch(cmd, disableMouseCmd)
}

func (sp *Splash) UpdateSize(width, height int) {
	sp.width = width
	sp.height = height
}

func (sp *Splash) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd = sp.logo.Update(msg)
	switch msg.(type) {
	case BrokerDataLoaded:
		slog.Info("Broker data loaded")
		sp.data = true
		cmd = tea.Batch(sp.LoadCmds()...)
		return cmd
	case DelayCompleteMsg:
		slog.Info("Delay Set")
		sp.delay = true
	}
	return cmd
}

func (sp *Splash) View() string {
	var msg string

	return lipgloss.Place(
		sp.width,
		sp.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			"",
			"",
			"",
			sp.logo.View(),
			"",
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				sp.theme.TextError().Render(msg),
			),
		),
	)
}

func (sp *Splash) IsLoadingComplete() bool {
	return sp.data &&
		sp.delay && sp.wsClient.IsConnected()
}

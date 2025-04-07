package models

import (
	"log/slog"
	"p1/pkg/tui/theme"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Splash struct {
	data   bool
	delay  bool
	theme  *theme.Theme
	logo   *Logo
	width  int
	height int
	error  *VisibleError
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

func NewSplash(theme *theme.Theme) *Splash {
	cursor := NewCursor(theme, 700)
	return &Splash{
		theme: theme,
		logo:  NewLogo(theme, cursor),
	}
}

func (sp *Splash) Init() tea.Cmd {

	var cmd tea.Cmd

	cmd = tea.Batch(cmd, sp.logo.Init())

	cmd = tea.Batch(cmd, func() tea.Msg {
		return BrokerDataLoaded{}
	})

	cmd = tea.Batch(cmd, func() tea.Msg {
		return tea.DisableMouse()
	})

	return cmd
}

func (sp *Splash) UpdateSize(width, height int) {
	sp.width = width
	sp.height = height
}

func (sp *Splash) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case VisibleError:
		sp.error = &msg
	case BrokerDataLoaded:
		sp.data = true
		cmds = append(cmds, sp.LoadCmds()...)
	case DelayCompleteMsg:
		sp.delay = true
	case tea.WindowSizeMsg:
		slog.Info("SPLASH WINDOW SIZE", "width", msg.Width, "height", msg.Height)
		sp.width = msg.Width
		sp.height = msg.Height
	}
	return tea.Batch(cmds...)
}

func (sp *Splash) View() string {
	var msg string

	if sp.error != nil {
		msg = sp.error.Message
	}

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
		sp.delay
}

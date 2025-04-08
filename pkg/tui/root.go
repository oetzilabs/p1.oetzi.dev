package tui

import (
	"context"
	"log/slog"
	"time"

	"p1/pkg/client"
	"p1/pkg/messages"
	"p1/pkg/models"
	"p1/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page = int

const (
	splashPage page = iota
	dashboardPage
)

type model struct {
	switched  bool
	renderer  *lipgloss.Renderer
	page      page
	dashboard *Dashboard
	context   context.Context
	theme     theme.Theme
	error     *models.VisibleError
	width     int
	height    int
	client    *client.Client
}

func NewModel(renderer *lipgloss.Renderer, client *client.Client) (tea.Model, error) {
	basicTheme := theme.BasicTheme(renderer, nil)

	dashboard := NewDashboard(&basicTheme)

	result := model{
		renderer:  renderer,
		context:   context.Background(),
		theme:     basicTheme,
		page:      dashboardPage,
		dashboard: dashboard,
		width:     0,
		height:    0,
		client:    client,
	}

	return result, nil
}

func (m *model) doSyncTick(client *client.Client) tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		slog.Info("Pulling")
		state := client.Pull()
		m.dashboard.SendMessageToTab("servers", messages.Message{Type: messages.TypeListServices, Payload: state.Servers})
		m.dashboard.SendMessageToTab("brokers", messages.Message{Type: messages.TypeListBrokers, Payload: state.Brokers})
		m.dashboard.SendMessageToTab("projects", messages.Message{Type: messages.TypeListProjects, Payload: state.Projects})
		return messages.SyncMsg(state)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.doSyncTick(m.client))
}

func (m model) SwitchPage(page page) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case messages.SyncMsg:
		return m, m.doSyncTick(m.client)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.error != nil {
				if m.page == splashPage {
					return m, tea.Quit
				}
				m.error = nil
				return m, nil
			}
		}
	}

	cmds = append(cmds, m.dashboard.Update(msg))

	if m.switched {
		m.switched = false
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string = m.dashboard.View()

	return m.renderer.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.theme.Base().
			MaxWidth(m.width).
			MaxHeight(m.height).
			Render(content),
	)
}

package tui

import (
	"context"
	"time"

	"p1/pkg/client"
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

type SyncMsg time.Time

type model struct {
	switched  bool
	renderer  *lipgloss.Renderer
	page      page
	dashboard *Dashboard
	splash    *models.Splash
	context   context.Context
	theme     theme.Theme
	error     *models.VisibleError
	width     int
	height    int
	client    *client.Client
}

func NewModel(renderer *lipgloss.Renderer) (tea.Model, error) {
	basicTheme := theme.BasicTheme(renderer, nil)

	splash := models.NewSplash(&basicTheme)
	dashboard := NewDashboard(&basicTheme)

	client := client.NewClient()

	result := model{
		renderer:  renderer,
		context:   context.Background(),
		theme:     basicTheme,
		page:      splashPage,
		splash:    splash,
		dashboard: dashboard,
		width:     0,
		height:    0,
		client:    client,
	}

	return result, nil
}

func doSyncTick(client *client.Client) tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		err := client.Sync()
		if err != nil {
			return tea.Batch(
				func() tea.Msg { return models.VisibleError{Message: err.Error()} },
				func() tea.Msg { return SyncMsg(t) },
			)()
		}
		return SyncMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.splash.Init(), doSyncTick(m.client))
}

func (m model) SwitchPage(page page) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case SyncMsg:
		return m, doSyncTick(m.client)
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

	cmds = append(cmds, m.splash.Update(msg))
	cmds = append(cmds, m.dashboard.Update(msg))

	if m.splash.IsLoadingComplete() && m.page == splashPage {
		m = m.SwitchPage(dashboardPage)
		return m, tea.Batch(cmds...)
	}

	if m.switched {
		m.switched = false
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {

	var content string
	switch m.page {
	case splashPage:
		content = m.splash.View()
	case dashboardPage:
		content = m.dashboard.View()
	}

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

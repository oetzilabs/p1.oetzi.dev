package tui

import (
	"context"

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
	splash    *models.Splash
	context   context.Context
	theme     theme.Theme
	error     *models.VisibleError
	width     int
	height    int
}

func NewModel(renderer *lipgloss.Renderer) (tea.Model, error) {
	basicTheme := theme.BasicTheme(renderer, nil)

	splash := models.NewSplash(&basicTheme)
	dashboard := NewDashboard(&basicTheme)

	result := model{
		renderer:  renderer,
		context:   context.Background(),
		theme:     basicTheme,
		page:      splashPage,
		splash:    splash,
		dashboard: dashboard,
		width:     0,
		height:    0,
	}

	return result, nil
}

func (m model) Init() tea.Cmd {
	return m.splash.Init()
}

func (m model) SwitchPage(page page) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case models.VisibleError:
		m.error = &msg
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
	if m.error != nil {
		return m.theme.TextError().Render("Error: " + m.error.Message)
	}

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

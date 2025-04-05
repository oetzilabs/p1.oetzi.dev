package tui

import (
	"context"

	"p1/pkg/models"
	"p1/pkg/tui/theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page = int

const (
	splashPage page = iota
	dashboardPage
)

type model struct {
	ready           bool
	switched        bool
	hasScroll       bool
	renderer        *lipgloss.Renderer
	page            page
	dashboard       *Dashboard
	splash          *models.Splash
	context         context.Context
	viewportWidth   int
	viewportHeight  int
	widthContainer  int
	heightContainer int
	widthContent    int
	heightContent   int
	viewport        viewport.Model
	theme           theme.Theme
	error           *models.VisibleError
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
		m.viewportWidth = msg.Width
		m.viewportHeight = msg.Height

		m.widthContainer = msg.Width
		m.heightContainer = msg.Height - 2 // Leave some margin

		m.widthContent = m.widthContainer - 4
		m = m.updateViewport()
		m.splash.UpdateSize(m.viewportWidth, m.viewportHeight)
		m.dashboard.UpdateSize(m.viewportWidth, m.viewportHeight)
		m.dashboard.sidebar.UpdateHeight(m.viewportHeight)
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

	var cmd tea.Cmd = tea.Batch(m.splash.Update(msg), m.dashboard.Update(msg))

	if m.splash.IsLoadingComplete() && m.page == splashPage {
		m = m.SwitchPage(dashboardPage)
		return m, cmd
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.viewport.SetContent(m.getContent())
	m.viewport, cmd = m.viewport.Update(msg)
	if m.switched {
		m = m.updateViewport()
		m.switched = false
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.error != nil {
		return m.theme.TextError().Render("Error: " + m.error.Message)
	}

	switch m.page {
	case splashPage:
		return m.splash.View()
	case dashboardPage:
		return m.dashboard.View()
	default:
		content := m.viewport.View()
		var view string
		if m.hasScroll {
			view = lipgloss.JoinHorizontal(
				lipgloss.Top,
				content,
				m.theme.Base().Width(1).Render(),
				m.getScrollbar(),
			)
		} else {
			view = m.getContent()
		}

		height := m.heightContainer

		child := lipgloss.JoinVertical(
			lipgloss.Left,
			m.theme.Base().
				Width(m.widthContainer).
				Height(height).
				Padding(0, 0).
				Render(view),
		)

		return m.renderer.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			m.theme.Base().
				MaxWidth(m.widthContainer).
				MaxHeight(m.heightContainer).
				Render(child),
		)
	}
}

func (m model) getContent() string {
	page := "unknown"
	switch m.page {
	case dashboardPage:
		page = m.dashboard.View()
	}
	return page
}

func (m model) getScrollbar() string {
	y := m.viewport.YOffset
	vh := m.viewport.Height
	ch := lipgloss.Height(m.getContent())
	if vh >= ch {
		return ""
	}

	height := (vh * vh) / ch
	maxScroll := ch - vh
	nYP := 1.0 - (float64(y) / float64(maxScroll))
	if nYP <= 0 {
		nYP = 1
	} else if nYP >= 1 {
		nYP = 0
	}

	bar := m.theme.Base().
		Height(height).
		Width(1).
		Background(m.theme.Accent()).
		Render()

	style := m.theme.Base().Width(1).Height(vh)

	return style.Render(lipgloss.PlaceVertical(vh, lipgloss.Position(nYP), bar))
}

var modifiedKeyMap = viewport.KeyMap{
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "½ page down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
}

func (m model) updateViewport() model {
	width := m.widthContainer

	verticalMarginHeight := 2

	m.heightContent = m.heightContainer - verticalMarginHeight

	if !m.ready {
		m.viewport = viewport.New(width, m.heightContent)
		m.viewport.HighPerformanceRendering = false
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = m.heightContent
		m.viewport.GotoTop()
	}

	m.hasScroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()

	if m.hasScroll {
		m.widthContent = m.widthContainer - 4
	} else {
		m.widthContent = m.widthContainer - 0
	}

	m.viewport.KeyMap = modifiedKeyMap

	return m
}

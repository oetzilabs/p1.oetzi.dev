package tui

import (
	"context"
	"errors"

	"p1/pkg/client"
	"p1/pkg/models"
	"p1/pkg/tui/theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page = int
type size = int

const (
	splashPage page = iota
	dashboardPage
)

const (
	undersized size = iota
	small
	medium
	large
)

type model struct {
	ready           bool
	switched        bool
	hasScroll       bool
	renderer        *lipgloss.Renderer
	page            page
	dashboard       *Dashboard
	cursor          Cursor
	splash          Splash
	context         context.Context
	viewportWidth   int
	viewportHeight  int
	widthContainer  int
	heightContainer int
	widthContent    int
	heightContent   int
	size            size
	viewport        viewport.Model
	theme           theme.Theme
	error           *VisibleError
	wsClient        *WebSocketClient // Will be implemented later
	wsConnected     bool
	client          *client.Client
	updateSub       chan models.Update
}

type VisibleError struct {
	message string
}

// WebSocketUpdateMsg represents a message received from the websocket
type WebSocketUpdateMsg struct {
	Type string
	Data interface{}
}

func NewModel(
	renderer *lipgloss.Renderer,
	wsURL string,
	command []string,
) (tea.Model, error) {
	if wsURL == "" {
		return nil, errors.New("WEBSOCKET_URL is not set")
	}

	client := client.NewClient(wsURL)

	go client.Start(context.Background())

	result := model{
		client:    client,
		updateSub: client.Subscribe(),
		context:   context.Background(),
		page:      splashPage,
		renderer:  renderer,
		theme:     theme.BasicTheme(renderer, nil),
	}

	return result, nil
}

func (m model) Init() tea.Cmd {
	return m.SplashInit()
}

func (m model) SwitchPage(page page) model {
	m.page = page
	m.switched = true
	return m
}

func (m model) InitialDataLoaded() (model, tea.Cmd) {
	return m.DashboardSwitch()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case VisibleError:
		m.error = &msg

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.viewportHeight = msg.Height

		switch {
		case m.viewportWidth < 20 || m.viewportHeight < 10:
			m.size = undersized
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 50:
			m.size = small
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 75:
			m.size = medium
			m.widthContainer = 50
			m.heightContainer = msg.Height - 2 // Leave some margin
		default:
			m.size = large
			m.widthContainer = 75
			m.heightContainer = msg.Height - 2 // Leave some margin
		}

		m.widthContent = m.widthContainer - 4
		m = m.updateViewport()
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
		case "ctrl+c":
			if m.wsClient != nil {
				m.wsClient.Disconnect()
			}
			return m, tea.Quit
		}
	case CursorTickMsg:
		m, cmd := m.CursorUpdate(msg)
		m.viewport.SetContent(m.getContent())
		return m, cmd
	case WebSocketUpdateMsg:
		// Just update the UI state based on the message
		// All connection/retry logic is handled by the client
	}

	var cmd tea.Cmd
	switch m.page {
	case splashPage:
		m, cmd = m.SplashUpdate(msg)
	case dashboardPage:
		m, cmd = m.DashboardUpdate(msg)
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

	// Always keep listening for updates
	cmds = append(cmds, listenForUpdates(m.updateSub))

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.size == undersized {
		return m.ResizeView()
	}

	switch m.page {
	case splashPage:
		return m.SplashView()
	case dashboardPage:
		return m.DashboardView()
	default:
		content := m.viewport.View()

		var view string
		if m.hasScroll {
			view = lipgloss.JoinHorizontal(
				lipgloss.Top,
				content,
				m.theme.Base().Width(1).Render(), // space between content and scrollbar
				m.getScrollbar(),
			)
		} else {
			view = m.getContent()
		}

		height := m.heightContainer
		// height -= lipgloss.Height(header)
		// height -= lipgloss.Height(breadcrumbs)
		// height -= lipgloss.Height(footer)

		child := lipgloss.JoinVertical(
			lipgloss.Left,
			// header,
			// breadcrumbs,
			m.theme.Base().
				Width(m.widthContainer).
				Height(height).
				Padding(0, 0).
				Render(view),
			// footer,
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
		page = m.DashboardView()
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

// Add a command to listen for updates
func listenForUpdates(sub chan models.Update) tea.Cmd {
	return func() tea.Msg {
		update := <-sub
		return WebSocketUpdateMsg{
			Type: update.Type,
			Data: update.Data,
		}
	}
}

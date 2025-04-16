package tui

import (
	"context"

	"p1/pkg/interfaces"
	"p1/pkg/menu"
	"p1/pkg/models"
	"p1/pkg/screens"
	"p1/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	renderer *lipgloss.Renderer
	context  context.Context
	theme    theme.Theme
	width    int
	height   int
	footer   *models.Footer
	menu     *menu.Menu
}

func NewModel(renderer *lipgloss.Renderer) tea.Model {
	basicTheme := theme.BasicTheme(renderer, nil)

	footerCommands := []*interfaces.FooterCommand{}

	footer := models.NewFooter(&basicTheme, footerCommands)
	footer.ResetCommands()

	default_menu := menu.NewMenu().
		AddItem(menu.NewMenuItem("projects", "Projects", screens.NewProjectsScreen(renderer))).
		AddItem(menu.NewMenuItem("brokers", "Brokers", screens.NewBrokersScreen(renderer)))

	result := model{
		renderer: renderer,
		context:  context.Background(),
		theme:    basicTheme,
		width:    0,
		height:   0,
		footer:   footer,
		menu:     default_menu,
	}

	return result
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return tea.DisableMouse() }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	parentMsg := msg
	cmds = append(cmds, m.menu.Update(parentMsg))
	switch msg := msg.(type) {
	// case messages.RerenderMessage:
	// 	slog.Info("Rerendering", "key", msg.Key, "value", msg.Value, "oldValue", msg.OldValue)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		cmds = append(cmds, tea.Cmd(func() tea.Msg {
			return tea.Msg(models.InternalWindowSizeMsg{
				Width:        msg.Width,
				Height:       msg.Height,
				MenuWidth:    m.menu.GetWidth(),
				FooterHeight: lipgloss.Height(m.footer.View()),
			})
		}))
	}

	m.footer.SetCommands(m.menu.GetCommands())

	cmds = append(cmds, m.footer.Update(parentMsg))

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	menu := m.menu.View()
	screen := m.menu.Screen()
	footerContent := m.footer.View()

	mainContent := lipgloss.JoinVertical(
		lipgloss.Top,
		screen,
		footerContent,
	)

	return m.renderer.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.theme.Base().
			MaxWidth(m.width).
			MaxHeight(m.height).
			Render(lipgloss.JoinHorizontal(lipgloss.Left, menu, mainContent)),
	)
}

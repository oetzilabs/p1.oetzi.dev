package screens

import (
	"p1/pkg/interfaces"
	"p1/pkg/models"
	"p1/pkg/tui/theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen struct {
	ready     bool
	Content   interfaces.ScreenContent
	viewport  viewport.Model
	menuWidth int
	height    int
	theme     theme.Theme
}

func NewScreen(renderer *lipgloss.Renderer, content interfaces.ScreenContent) *Screen {
	return &Screen{
		Content: content,
		ready:   false,
		theme:   theme.BasicTheme(renderer, nil),
	}
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
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	),
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}
	parentMsg := msg
	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		s.height = msg.Height
		if !s.ready {
			s.viewport = viewport.New(msg.Width-msg.MenuWidth, msg.Height-msg.FooterHeight)
			s.viewport.HighPerformanceRendering = false
			s.viewport.KeyMap = modifiedKeyMap
			s.viewport.Style = s.viewport.Style.
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				Height(msg.Height - msg.FooterHeight).
				Width(msg.Width - msg.MenuWidth).
				MaxHeight(msg.Height - msg.FooterHeight)
			s.ready = true
		} else {
			s.viewport.Width = max(0, msg.Width - msg.MenuWidth)
			s.viewport.Height = max(0, msg.Height - msg.FooterHeight)
		}
	}
	cmds = append(cmds, s.Content.Update(parentMsg))
	s.viewport.SetContent(s.Content.View())
	return tea.Batch(cmds...)
}

func (s *Screen) View() string {
	content := s.viewport.View()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		s.getScrollbar(content),
	)
}

func (s *Screen) getScrollbar(content string) string {
	y := s.viewport.YOffset
	vh := s.viewport.Height
	ch := lipgloss.Height(content)
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

	bar := s.theme.Base().
		Height(height).
		Width(1).
		Background(s.theme.Accent()).
		Render()

	style := s.theme.Base().Width(1).Height(vh)

	return style.Render(lipgloss.PlaceVertical(vh, lipgloss.Position(nYP), bar))
}

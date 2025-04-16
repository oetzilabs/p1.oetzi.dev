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
	ready        bool
	Content      interfaces.ScreenContent
	viewport     viewport.Model
	menuWidth    int
	footerHeight int
	width        int
	height       int
	theme        theme.Theme
	focused      bool
	Commands     []*interfaces.FooterCommand
}

func (s *Screen) IsFocused() bool {
	return s.focused
}

func (s *Screen) SetFocused(focused bool) *Screen {
	s.focused = focused
	return s
}

func NewScreen(renderer *lipgloss.Renderer, content interfaces.ScreenContent, commands ...*interfaces.FooterCommand) *Screen {
	cmds := []*interfaces.FooterCommand{}
	cmds = append(cmds, commands...)
	return &Screen{
		Content:  content,
		ready:    false,
		theme:    theme.BasicTheme(renderer, nil),
		focused:  false,
		Commands: cmds,
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
		s.menuWidth = msg.MenuWidth
		s.width = msg.Width
		s.footerHeight = msg.FooterHeight
		if !s.ready {
			s.viewport = viewport.New(msg.Width-msg.MenuWidth-2, msg.Height-msg.FooterHeight)
			s.viewport.HighPerformanceRendering = false
			s.viewport.KeyMap = modifiedKeyMap
			s.viewport.Style = s.viewport.Style.
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				PaddingBottom(1)
			s.ready = true
		} else {
			s.viewport.Width = max(0, msg.Width-msg.MenuWidth-2)
			s.viewport.Height = max(0, msg.Height-msg.FooterHeight)
			s.viewport.Style = s.viewport.Style.MaxHeight(msg.Height - msg.FooterHeight)
		}
	}

	s.viewport.KeyMap = modifiedKeyMap
	s.viewport.Width = max(0, s.width-s.menuWidth-2)
	s.viewport.Height = max(0, s.height-s.footerHeight)
	s.viewport.Style = s.viewport.Style.MaxHeight(max(0, s.height-s.footerHeight))
	cmds = append(cmds, s.Content.Update(parentMsg))
	s.viewport.SetContent(s.Content.View())

	if s.focused {
		vpm, cmd := s.viewport.Update(parentMsg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		s.viewport = vpm
	}

	return tea.Batch(cmds...)
}

func (s *Screen) View() string {
	viewportView := s.viewport.View()
	content := s.Content.View()
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		viewportView,
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

func (s *Screen) Display() string {
	if s.Content != nil {
		return s.Content.Display()
	}
	return ""
}

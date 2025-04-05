package models

import (
	"p1/pkg/tui/theme"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Cursor struct {
	visible    bool
	theme      *theme.Theme
	tickMillis int
}

type CursorTickMsg struct{}

func NewCursor(theme *theme.Theme, tickMillis int) *Cursor {
	return &Cursor{
		theme:      theme,
		visible:    true,
		tickMillis: tickMillis,
	}
}

func (c *Cursor) Init() tea.Cmd {
	return tea.Every(time.Millisecond*time.Duration(c.tickMillis), func(t time.Time) tea.Msg {
		return CursorTickMsg{}
	})
}

func (c *Cursor) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case CursorTickMsg:
		c.visible = !c.visible
		return tea.Every(time.Millisecond*time.Duration(c.tickMillis), func(t time.Time) tea.Msg {
			return CursorTickMsg{}
		})
	}
	return nil
}

func (c *Cursor) View() string {
	if c.visible {
		return c.theme.Base().Background(c.theme.Highlight()).Render(" ")
	} else {
		return c.theme.Base().Render(" ")
	}
}

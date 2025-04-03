package models

import (
	"log/slog"
	"p1/pkg/tui/theme"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Cursor struct {
	visible bool
	theme   *theme.Theme
}

type CursorTickMsg struct{}

func NewCursor(theme *theme.Theme) *Cursor {
	return &Cursor{
		theme:   theme,
		visible: true,
	}
}

func (c *Cursor) Init() tea.Cmd {
	slog.Info("Initializing Cursor")
	return tea.Every(time.Millisecond*700, func(t time.Time) tea.Msg {
		return CursorTickMsg{}
	})
}

func (c *Cursor) Update(msg tea.Msg) tea.Cmd {
	// slog.Info("Updating Cursor")
	switch msg.(type) {
	case CursorTickMsg:
		c.visible = !c.visible
		return tea.Every(time.Millisecond*700, func(t time.Time) tea.Msg {
			slog.Info("Toggling Cursor")
			return CursorTickMsg{}
		})
	}
	return nil
}

func (c *Cursor) View() string {
	// slog.Info("Rendering Cursor")
	if c.visible {
		return c.theme.Base().Background(c.theme.Highlight()).Render(" ")
	} else {
		return c.theme.Base().Render(" ")
	}
}

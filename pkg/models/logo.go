package models

import (
	"p1/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
)

type Logo struct {
	theme  *theme.Theme
	cursor *Cursor
}

func NewLogo(theme *theme.Theme, cursor *Cursor) *Logo {
	return &Logo{
		theme:  theme,
		cursor: cursor,
	}
}

func (l *Logo) View() string {
	// slog.Info("Rendering Logo")
	return l.theme.TextAccent().Bold(true).Render("p1.oetzi.dev ") + l.cursor.View()
}

func (l *Logo) Update(msg tea.Msg) tea.Cmd {
	// slog.Info("Updating Logo")
	return l.cursor.Update(msg)
}

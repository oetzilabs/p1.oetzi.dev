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

func (l *Logo) Init() tea.Cmd {
	return l.cursor.Init()
}

func (l *Logo) View() string {
	return l.theme.TextAccent().Bold(true).Render("p1.oetzi.dev ") + l.cursor.View()
}

func (l *Logo) Update(msg tea.Msg) tea.Cmd {
	return l.cursor.Update(msg)
}

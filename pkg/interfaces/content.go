package tui

import tea "github.com/charmbracelet/bubbletea"

type Content interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
	Display() string
}

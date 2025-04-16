package interfaces

import tea "github.com/charmbracelet/bubbletea"

type FooterCommand struct {
	Key   string
	Value string
}

type Content interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
	Display() string
	Commands() []*FooterCommand
}

type ScreenContent interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
	Display() string
}

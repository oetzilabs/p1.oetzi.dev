package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Confirmation struct {
	message   string
	confirmed bool
	onConfirm func()
	onCancel  func()
}

func NewConfirmation(message string, onConfirm, onCancel func()) *Confirmation {
	return &Confirmation{
		message:   message,
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
}

func (c *Confirmation) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			c.confirmed = true
			c.onConfirm()
		case "n":
			c.confirmed = false
			c.onCancel()
		}
	}
	return nil
}

func (c *Confirmation) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		c.message,
		"(y) Yes  |  (n) No",
	)
}

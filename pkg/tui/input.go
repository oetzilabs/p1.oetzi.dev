package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	Label string
	Value string
}

func NewInput(label string, value string) *Input {
	return &Input{
		Label: label,
		Value: value,
	}
}

func (i *Input) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return tea.Quit
		}
	}
	return nil
}

func (i *Input) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, i.Label, i.Value)
}

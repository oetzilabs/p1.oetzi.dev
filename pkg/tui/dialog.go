package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dialog struct {
	message   string
	inputs    []Input
	onConfirm func(inputs ...string)
	onCancel  func()
}

func NewDialog(message string, onConfirm func(inputs ...string), onCancel func(), inputs ...Input) *Dialog {
	return &Dialog{
		message:   message,
		inputs:    inputs,
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
}

func (d *Dialog) Update(msg tea.Msg) tea.Cmd {
	inputs := []string{}
	for _, input := range d.inputs {
		inputs = append(inputs, input.Value)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			d.onConfirm(inputs...)
		case "n":
			d.onCancel()
		}
	}
	return nil
}

func (d *Dialog) View() string {
	inputs := []string{}
	for _, input := range d.inputs {
		inputs = append(inputs, input.View())
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		d.message,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			inputs...,
		),
		"(y) Yes  |  (n) No",
	)
}

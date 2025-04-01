package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputOptions struct {
	Focused     bool
	Disabled    bool
	Placeholder string
}

type Input struct {
	Label   string
	Value   string
	Options InputOptions
}

// DefaultInputOptions returns the default options for an input
func DefaultInputOptions() InputOptions {
	return InputOptions{
		Focused:     false,
		Disabled:    false,
		Placeholder: "",
	}
}

// NewInput creates a new input with optional options
// If options is nil, default options will be used
func NewInput(label string, value string, options *InputOptions) *Input {
	opts := DefaultInputOptions()
	if options != nil {
		opts = *options
	}
	return &Input{
		Label:   label,
		Value:   value,
		Options: opts,
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

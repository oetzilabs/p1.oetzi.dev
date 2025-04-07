package models

import (
	"fmt"
	"p1/pkg/interfaces"
	"p1/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Footer struct {
	Commands []interfaces.FooterCommand
	theme    *theme.Theme
	error    *VisibleError
	width    int
	helper   string
}

type FooterUpdate struct {
	Content  string
	Commands []interfaces.FooterCommand
}

var (
	BaseCommands = []interfaces.FooterCommand{
		{Key: "q", Value: "quit"},
		{Key: "ctrl+k", Value: "Focus Sidebar"},
	}
)

func NewFooter(theme *theme.Theme, commands []interfaces.FooterCommand) *Footer {
	return &Footer{
		theme:    theme,
		Commands: commands,
		width:    0,
		helper:   "",
	}
}

func (f *Footer) UpdateWidth(width int) {
	f.width = width
}

func (f *Footer) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case VisibleError:
		f.error = &msg
	case FooterUpdate:
		f.helper = msg.Content
		f.Commands = msg.Commands
	}
	return nil
}

func (f *Footer) ResetCommands() {
	f.Commands = BaseCommands
}

func (f *Footer) View() string {
	bold := f.theme.TextAccent().Background(lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FFFFFF"}).Bold(true).Render
	base := f.theme.TextAccent().Background(lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FFFFFF"}).Render

	table := f.theme.Base().
		Width(f.width - 2).
		Background(lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FFFFFF"}).
		Padding(1).
		PaddingLeft(2).
		PaddingRight(2).
		Align(lipgloss.Right)

	lines := []string{}

	content := lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FFFFFF"}).
		Foreground(lipgloss.Color("#777777")).Render(f.helper + " ")

	if f.error != nil {
		errorString := fmt.Sprintf("Error: %s ", f.error.Message)
		content = lipgloss.JoinHorizontal(
			lipgloss.Left,
			f.theme.TextError().Render(errorString),
			content,
		)
	}

	lines = append(lines, content)

	mergedCommands := BaseCommands
	mergedCommands = append(mergedCommands, f.Commands...)

	// Add other commands
	commands := []string{}
	for cmdIndex, cmd := range mergedCommands {
		spacer := ""
		if cmdIndex < len(f.Commands)-1 {
			spacer = "|"
		}
		commands = append(commands, bold(" "+cmd.Key+" ")+base(cmd.Value+" ")+base(spacer))
	}

	lines = append(lines, commands...)

	return table.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Right,
			lines...,
		),
	)
}

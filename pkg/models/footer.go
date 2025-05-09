package models

import (
	"fmt"
	"p1/pkg/interfaces"
	"p1/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Footer struct {
	Commands []*interfaces.FooterCommand
	theme    *theme.Theme
	error    *VisibleError
	width    int
	helper   string
}

type FooterUpdate struct {
	Content  string
	Commands []*interfaces.FooterCommand
}

var (
	BaseCommands = []*interfaces.FooterCommand{
		{Key: "q", Value: "Quit"},
		{Key: "ctrl+k", Value: "Focus Sidebar"},
	}
)

func NewFooter(theme *theme.Theme, commands []*interfaces.FooterCommand) *Footer {
	return &Footer{
		theme:    theme,
		Commands: commands,
		width:    0,
		helper:   "",
	}
}

func (f *Footer) SetCommands(cmds []*interfaces.FooterCommand) {
	f.ResetCommands()
	f.Commands = append(f.Commands, cmds...)
}

func (f *Footer) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case VisibleError:
		f.error = &msg
	case InternalWindowSizeMsg:
		f.width = msg.Width - msg.MenuWidth
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
		Width(f.width - 1).
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
	for _, cmd := range f.Commands {
		// only add non existing keys
		found := false
		for _, cmd2 := range BaseCommands {
			if cmd.Key == cmd2.Key {
				found = true
				break
			}
		}
		if !found {
			mergedCommands = append(mergedCommands, cmd)
		}
	}

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

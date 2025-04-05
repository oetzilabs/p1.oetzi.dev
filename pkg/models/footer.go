package models

import (
	"p1/pkg/tui/theme"
	"p1/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Footer struct {
	Commands []FooterCommand
	theme    *theme.Theme
	error    *VisibleError
	width    int
	helper   string
}

type FooterCommand struct {
	Key   string
	Value string
}

type FooterUpdateHelperMsg struct {
	Content string
}

type FooterUpdateCommandsMsg struct {
	Commands []FooterCommand
}

func NewFooter(theme *theme.Theme, commands []FooterCommand) *Footer {
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
	case FooterUpdateHelperMsg:
		f.helper = msg.Content
	case FooterUpdateCommandsMsg:
		f.Commands = msg.Commands
	}
	return nil
}

func (f *Footer) ResetCommands() {
	f.Commands = []FooterCommand{
		{Key: "q", Value: "quit"},
		{Key: "ctrl+k", Value: "Focus Sidebar"},
	}
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

	var content string
	if f.error != nil {
		hint := "esc"

		// Calculate maximum width for error message to ensure it fits
		maxErrorWidth := f.width - lipgloss.Width(hint) - 6

		// Handle wrapping for long error messages
		errorMsg := f.error.Message
		if lipgloss.Width(errorMsg) > maxErrorWidth {
			// Split into multiple lines
			errorMsg = utils.WordWrap(errorMsg, maxErrorWidth)
		}

		msg := f.theme.PanelError().Padding(0, 1).Render(errorMsg)

		// Calculate remaining space after rendering the message
		space := max(0, f.width-lipgloss.Width(msg)-lipgloss.Width(hint)-2)

		height := lipgloss.Height(msg)

		content = lipgloss.JoinVertical(
			lipgloss.Right,
			msg,
			f.theme.PanelError().Width(space).Height(height).Render(),
			f.theme.PanelError().Bold(true).Padding(0, 1).Height(height).Render(hint),
		)
	} else {
		content = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FFFFFF"}).
			Foreground(lipgloss.Color("#777777")).Render(f.helper + " ")
	}

	lines = append(lines, content)

	// Add other commands
	commands := []string{}
	for cmdIndex, cmd := range f.Commands {
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

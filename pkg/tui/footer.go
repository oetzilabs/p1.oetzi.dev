package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// wordWrap breaks a string into multiple lines to fit within maxWidth
func wordWrap(text string, maxWidth int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	lines := []string{}
	currentLine := words[0]

	for _, word := range words[1:] {
		// Check if adding this word would exceed the width
		testLine := currentLine + " " + word
		if lipgloss.Width(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			// Line would be too long, start a new line
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	// Add the last line
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

func (m model) FooterView() string {
	bold := m.theme.TextAccent().Bold(true).Render
	base := m.theme.Base().Render

	table := m.theme.Base().
		Width(m.widthContainer).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.theme.Border()).
		PaddingBottom(1).
		Align(lipgloss.Center)

	// Add other commands
	commands := []string{}
	for _, cmd := range m.footer.Commands {
		commands = append(commands, bold(" "+cmd.Key+" ")+base(cmd.Value+"  "))
	}

	lines := []string{}
	// if m.page == dashboardPage {
	// 	lines = append(lines, bold("r")+regionSelector)
	// 	lines = append(lines, base("  "))
	// }
	lines = append(lines, commands...)

	var content string
	// if m.error != nil {
	// 	hint := "esc"

	// 	// Calculate maximum width for error message to ensure it fits
	// 	maxErrorWidth := m.widthContent - lipgloss.Width(hint) - 6

	// 	// Handle wrapping for long error messages
	// 	errorMsg := m.error.message
	// 	if lipgloss.Width(errorMsg) > maxErrorWidth {
	// 		// Split into multiple lines
	// 		errorMsg = wordWrap(errorMsg, maxErrorWidth)
	// 	}

	// 	msg := m.theme.PanelError().Padding(0, 1).Render(errorMsg)

	// 	// Calculate remaining space after rendering the message
	// 	space := m.widthContent - lipgloss.Width(msg) - lipgloss.Width(hint) - 2
	// 	if space < 0 {
	// 		space = 0
	// 	}

	// 	height := lipgloss.Height(msg)

	// 	content = lipgloss.JoinHorizontal(
	// 		lipgloss.Top,
	// 		msg,
	// 		m.theme.PanelError().Width(space).Height(height).Render(),
	// 		m.theme.PanelError().Bold(true).Padding(0, 1).Height(height).Render(hint),
	// 	)
	// } else {
	// 	content = "test test test"
	// }

	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		content,
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				lines...,
			),
		))
}

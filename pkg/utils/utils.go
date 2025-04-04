package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func WordWrap(text string, maxWidth int) string {
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

package utils

import (
	"crypto/rand"
	"math/big"
	"p1/pkg/tabs"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Tabs = []tabs.Tab

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

func GenerateLoremIpsum(length int) string {
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	var words []string
	for i := 0; i < length; i++ {
		// random number between 1 and 4
		randomnumber, err := rand.Int(rand.Reader, big.NewInt(4))
		if err != nil {
			continue
		}
		words = append(words, strings.Repeat(lorem, int(randomnumber.Int64())+1))
	}
	return strings.Join(words, "\n")
}

func FilterTabs(tabs Tabs, search string) Tabs {
	if search == "" {
		return tabs
	}
	var tabsToRender Tabs
	for _, tab := range tabs {
		if tab.IgnoreSearch {
			tabsToRender = append(tabsToRender, tab)
			continue
		}
		display := strings.ToLower(tab.Display())
		if len(display) == 0 {
			continue
		}

		search := strings.ToLower(search)
		if strings.Contains(display, search) {
			tabsToRender = append(tabsToRender, tab)
		}
	}
	return tabsToRender
}

func FilterTabsGroup(tabs Tabs, group tabs.TabPosition) Tabs {
	var tabsToRender Tabs
	for _, tab := range tabs {
		if tab.Group == group {
			tabsToRender = append(tabsToRender, tab)
		}
	}
	return tabsToRender
}

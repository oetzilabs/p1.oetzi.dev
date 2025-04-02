package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"

	tabs "p1/pkg/tabs"
)

type Sidebar struct {
	tabs         []tabs.Tab
	activeTab    int
	focused      bool
	search       textinput.Model
	tabsToRender []tabs.Tab
}

type Tabs = []tabs.Tab

func NewSidebar(tabs ...tabs.Tab) *Sidebar {
	ti := textinput.New()
	ti.Placeholder = "Search"

	return &Sidebar{
		tabs:         tabs,
		activeTab:    0,
		focused:      true,
		search:       ti,
		tabsToRender: tabs,
	}
}

func filterTabs(tabs Tabs, search string) Tabs {
	var tabsToRender Tabs
	for _, tab := range tabs {
		display := strings.ToLower(tab.Display())
		search := strings.ToLower(search)
		if strings.Contains(display, search) {
			tabsToRender = append(tabsToRender, tab)
		}
	}
	return tabsToRender
}

func filterTabsGroup(tabs Tabs, group tabs.TabGroup) Tabs {
	var tabsToRender Tabs
	for _, tab := range tabs {
		if tab.Group == group {
			tabsToRender = append(tabsToRender, tab)
		}
	}
	return tabsToRender
}

func (s *Sidebar) ViewSelectedTabContent() string {

	if len(s.tabsToRender) == 0 {
		return "No tabs found"
	}

	return s.tabsToRender[s.activeTab].Content.View()
}

func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	tiM, cmd := s.search.Update(msg)
	s.search = tiM

	s.tabsToRender = filterTabs(
		s.tabs,
		s.search.Value(),
	)

	if len(s.tabsToRender) == 0 {
		s.activeTab = 0
		return cmd
	}

	tab := s.tabsToRender[s.activeTab]
	tabCmd := tab.Update(msg)
	cmd = tea.Batch(cmd, tabCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			s.search.Focus()
		case "esc":
			if s.search.Focused() {
				s.search.Blur()
			} else {
				s.focused = false
			}
		case "j", "down":
			// Move down in sidebar
			if !s.search.Focused() && s.activeTab < len(s.tabsToRender)-1 && s.focused {
				s.activeTab++
			}
		case "k", "up":
			// Move up in sidebar
			if !s.search.Focused() && s.activeTab > 0 && s.focused {
				s.activeTab--
			}

		case "enter", "tab":
			// unfocus sidebar
			if !s.search.Focused() {
				s.focused = false
			}
		case "ctrl+k":
			// Return to sidebar
			if !s.search.Focused() {
				s.focused = true
			}
		}
	}
	return cmd
}

// Sidebar's View method
func (s *Sidebar) View() string {
	paddedStyle := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingTop(1).AlignVertical(lipgloss.Top)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	inactiveSelectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Bold(true)
	sidebarStyle := lipgloss.NewStyle().PaddingLeft(2)

	var sidebar string

	sidebar += s.search.View() + "\n\n"

	availableWidth := 30

	mainTabs := filterTabsGroup(s.tabsToRender, tabs.TabGroupsMain)
	bottomTabs := filterTabsGroup(s.tabsToRender, tabs.TabGroupsBottom)

	for _, tab := range mainTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()

		if len(display) > availableWidth {
			display = display[:availableWidth] + "..."
		}

		if s.tabsToRender[s.activeTab].ID == tab.ID {
			if s.focused {
				sidebar += selectedStyle.Render("▶ "+display) + "\n"
			} else {
				sidebar += inactiveSelectedStyle.Render("▶ "+display) + "\n"
			}
		} else {
			sidebar += sidebarStyle.Render(display) + "\n"
		}
	}

	sidebar += "__FILLER_VERTICAL__"

	for _, tab := range bottomTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()

		if len(display) > availableWidth {
			display = display[:availableWidth] + "..."
		}

		if s.tabsToRender[s.activeTab].ID == tab.ID {
			if s.focused {
				sidebar += selectedStyle.Render("▶ "+display) + "\n"
			} else {
				sidebar += inactiveSelectedStyle.Render("▶ "+display) + "\n"
			}
		} else {
			sidebar += sidebarStyle.Render(display) + "\n"
		}
	}

	return paddedStyle.Render(lipgloss.JoinVertical(lipgloss.Top, sidebar))
}

func (s *Sidebar) SidebarView(m model) string {
	// Replace filler string with spaces
	filler := "__FILLER__"
	replacement := " " // Single space
	original := s.View()
	// Create the replacement string with the correct number of spaces
	fillerWidth := max(0, 29-(lipgloss.Width(original)-lipgloss.Width(filler)))
	newString := strings.ReplaceAll(original, filler, strings.Repeat(replacement, fillerWidth))

	// Replace filler string with spaces, so that we fill the vertical space
	fillerVertical := "__FILLER_VERTICAL__"
	replacementVertical := "\n" // Single space
	fillerVerticalHeight := m.viewportHeight - lipgloss.Height(newString)

	newStringVertical := strings.ReplaceAll(newString, fillerVertical, strings.Repeat(replacementVertical, fillerVerticalHeight))

	return lipgloss.NewStyle().Background(lipgloss.Color("#111111")).Width(30).Height(m.viewportHeight).Render(newStringVertical)
}

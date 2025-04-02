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
	width        int
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
		width:        30,
	}
}

func filterTabs(tabs Tabs, search string) Tabs {
	var tabsToRender Tabs
	for _, tab := range tabs {
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

	content := s.tabsToRender[s.activeTab].Content
	if content == nil {
		return "No content found"
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
			if s.search.Focused() {
				s.search.Blur()
			}
			s.focused = true
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

	mainTabs := filterTabsGroup(s.tabsToRender, tabs.AlignTop)
	bottomTabs := filterTabsGroup(s.tabsToRender, tabs.AlignBottom)

	for _, tab := range mainTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()

		if len(display) > s.width {
			display = display[:s.width] + "..."
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

	sidebar += "__FILLER_VERTICAL__\n"

	for _, tab := range bottomTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()

		if len(display) > s.width {
			display = display[:s.width] + "..."
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
	original := s.View()
	sidebarWidth := m.dashboard.sidebar.width
	viewportHeight := m.viewportHeight

	lines := strings.Split(original, "\n")
	var processedLines []string

	for _, line := range lines {
		fillerIndex := strings.Index(line, "__FILLER__")
		if fillerIndex != -1 {
			// Extract content before and after the filler
			contentBeforeFiller := line[:fillerIndex]
			contentAfterFiller := line[fillerIndex+len("__FILLER__"):]

			// Calculate space needed for this specific line
			widthBeforeFiller := lipgloss.Width(contentBeforeFiller)
			spaceNeeded := max(0, sidebarWidth-widthBeforeFiller-lipgloss.Width(contentAfterFiller))
			horizontalSpace := strings.Repeat(" ", spaceNeeded)

			// Construct the new line
			line = contentBeforeFiller + horizontalSpace + contentAfterFiller
		}
		processedLines = append(processedLines, line)
	}

	// Join the lines back together
	content := strings.Join(processedLines, "\n")
	content = lipgloss.JoinHorizontal(lipgloss.Top, content)

	// Vertical spacing
	verticalFillerHeight := viewportHeight - lipgloss.Height(content)
	verticalSpace := strings.Repeat("\n", verticalFillerHeight)
	finalContent := strings.ReplaceAll(content, "__FILLER_VERTICAL__", verticalSpace)

	return lipgloss.NewStyle().
		Background(lipgloss.Color("#111111")).
		Width(s.width).
		Height(viewportHeight).
		Render(finalContent)
}

package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"

	"p1/pkg/interfaces"
	tabs "p1/pkg/tabs"
)

type Sidebar struct {
	tabs         []tabs.Tab
	activeTab    int
	focused      bool
	search       textinput.Model
	tabsToRender []tabs.Tab
	width        int
	height       int
}

type Tabs = []tabs.Tab

func NewSidebar(tabs ...tabs.Tab) *Sidebar {
	ti := textinput.New()
	ti.PromptStyle.MaxWidth(23)
	ti.PromptStyle.Width(23)
	ti.Prompt = "# "
	ti.Placeholder = "Search" + strings.Repeat(" ", 23-lipgloss.Width("Search"))

	return &Sidebar{
		tabs:         tabs,
		activeTab:    0,
		focused:      true,
		search:       ti,
		tabsToRender: tabs,
		width:        30,
		height:       len(tabs),
	}
}

func filterTabs(tabs Tabs, search string) Tabs {
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

func filterTabsGroup(tabs Tabs, group tabs.TabGroup) Tabs {
	var tabsToRender Tabs
	for _, tab := range tabs {
		if tab.Group == group {
			tabsToRender = append(tabsToRender, tab)
		}
	}
	return tabsToRender
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
	case tea.WindowSizeMsg:
		s.height = msg.Height
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
		case "ctrl+c":
			if s.search.Focused() {
				s.search.Blur()
			} else {
				return tea.Quit
			}
		case "q":
			if !s.search.Focused() {
				return tea.Quit
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
			s.focused = false
		case "ctrl+k":
			// Return to sidebar
			if !s.search.Focused() {
				s.focused = true
			}
		}
	}
	return cmd
}

func (s *Sidebar) formatTabEntry(display string, active bool) string {
	parts := strings.SplitN(display, " ", 2)
	title := parts[0]
	info := ""
	if len(parts) > 1 {
		info = parts[1]
	}

	sidebarWidth := s.width - 7 // Account for padding and cursor
	spaceWidth := max(1, sidebarWidth-lipgloss.Width(title)-lipgloss.Width(info))
	spacing := strings.Repeat(" ", spaceWidth)

	content := title + spacing + info
	if active {
		if s.focused {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				PaddingRight(2).
				Render("▶ " + content)
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Bold(true).
			PaddingRight(2).
			Render("▶ " + content)
	}
	return lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		Render(content)
}

func (s *Sidebar) View() string {
	paddedStyle := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingTop(1).AlignVertical(lipgloss.Top)
	var sidebar string

	sidebar += s.search.View() + "\n\n"

	mainTabs := filterTabsGroup(s.tabsToRender, tabs.AlignTop)
	bottomTabs := filterTabsGroup(s.tabsToRender, tabs.AlignBottom)

	for _, tab := range mainTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()
		isActive := s.tabsToRender[s.activeTab].ID == tab.ID
		sidebar += s.formatTabEntry(display, isActive) + "\n"
	}

	sidebar += "__FILLER_VERTICAL__\n"

	for _, tab := range bottomTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()
		isActive := s.tabsToRender[s.activeTab].ID == tab.ID
		sidebar += s.formatTabEntry(display, isActive) + "\n"
	}

	content := paddedStyle.Render(lipgloss.JoinVertical(lipgloss.Top, sidebar))

	content = lipgloss.JoinHorizontal(lipgloss.Top, content)

	// Vertical spacing
	verticalFillerHeight := max(0, s.height-lipgloss.Height(content))
	verticalSpace := strings.Repeat("\n", verticalFillerHeight)
	finalContent := strings.ReplaceAll(content, "__FILLER_VERTICAL__", verticalSpace)
	return lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Dark: "#111111", Light: "#EEEEEE"}).
		Width(s.width).
		Height(s.height).
		Render(finalContent)
}

func (s *Sidebar) SelectedTabContentView() string {
	if len(s.tabsToRender) == 0 {
		return "No tabs found"
	}

	content := s.tabsToRender[s.activeTab].Content
	if content == nil {
		return "No content found"
	}

	return content.View()
}
func (s *Sidebar) SelectedTabContentCommands() []interfaces.FooterCommand {
	if len(s.tabsToRender) == 0 {
		return []interfaces.FooterCommand{}
	}

	return s.tabsToRender[s.activeTab].Commands()
}

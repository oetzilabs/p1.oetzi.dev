package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
)

type Sidebar struct {
	tabs      []Tab
	activeTab int
	focused   bool
	search    textinput.Model
}

func NewSidebar(tabs ...Tab) *Sidebar {
	ti := textinput.New()
	ti.Placeholder = "Search"

	return &Sidebar{
		tabs:      tabs,
		activeTab: 0,
		focused:   true,
		search:    ti,
	}
}

func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {

	tab := s.tabs[s.activeTab]
	tabCmd := tab.Update(msg)
	tiM, cmd := s.search.Update(msg)
	s.search = tiM
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
		case "j":
			// Move down in sidebar
			if !s.search.Focused() && s.activeTab < len(s.tabs)-1 && s.focused {
				s.activeTab++
			}
		case "k":
			// Move up in sidebar
			if !s.search.Focused() && s.activeTab > 0 && s.focused {
				s.activeTab--
			}

		case "enter":
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
	paddedStyle := lipgloss.NewStyle().PaddingLeft(2)

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	inactiveSelectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Bold(true)
	sidebarStyle := lipgloss.NewStyle().PaddingLeft(1)

	var sidebar string

	sidebar += s.search.View() + "\n"

	for _, tab := range s.tabs {
		if tab.Hidden {
			continue
		}

		if s.tabs[s.activeTab].ID == tab.ID {
			if s.focused {
				sidebar += selectedStyle.Render("▶ "+tab.Display()) + "\n"
			} else {
				sidebar += inactiveSelectedStyle.Render("▶ "+tab.Display()) + "\n"
			}
		} else {
			sidebar += sidebarStyle.Render(tab.Display()) + "\n"
		}
	}

	return paddedStyle.Render(sidebar)
}

func (s *Sidebar) SidebarView(m model) string {

	borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderRight(true)

	sidebarBox := borderStyle.Width(20).Height(m.heightContainer).Render(s.View())

	return sidebarBox
}

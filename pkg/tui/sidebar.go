package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	tabs      []Tab
	activeTab int
	focused   bool
}

func NewSidebar(tabs ...Tab) *Sidebar {
	return &Sidebar{
		tabs:      tabs,
		activeTab: 0,
		focused:   true,
	}
}

func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for i := range s.tabs {
		// Update both the tab and its content
		if cmd := s.tabs[i].Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}

	}

	return tea.Batch(cmds...)
}

// Sidebar's View method
func (s *Sidebar) View() string {
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	inactiveSelectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Bold(true)
	sidebarStyle := lipgloss.NewStyle().PaddingLeft(1)

	var sidebar string
	for _, t := range s.tabs {
		if t.Hidden {
			continue
		}

		if s.tabs[s.activeTab].ID == t.ID {
			if s.focused {
				sidebar += selectedStyle.Render("▶ "+t.View()) + "\n"
			} else {
				sidebar += inactiveSelectedStyle.Render("▶ "+t.View()) + "\n"
			}
		} else {
			sidebar += sidebarStyle.Render(t.View()) + "\n"
		}
	}
	return sidebar
}

func (s *Sidebar) SidebarView(m model) string {
	borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderRight(true)

	sidebarContent := s.View()
	sidebarBox := borderStyle.Width(20).Height(m.heightContainer).Render(sidebarContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox)
}

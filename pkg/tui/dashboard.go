package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	sidebar *Sidebar
}

func NewDashboard() *Dashboard {
	projectsTab := NewTab("projects", NewProjectCollection())
	serversTab := NewTab("servers", NewServerCollection())
	brokersTab := NewTab("brokers", NewBrokerCollection())
	aboutTab := NewTab("about", NewAbout())

	return &Dashboard{
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			aboutTab,
		),
	}
}

func (m model) DashboardSwitch() (model, tea.Cmd) {
	m.dashboard = NewDashboard()
	m = m.SwitchPage(dashboardPage)
	return m, nil
}

func (m model) DashboardUpdate(msg tea.Msg) (model, tea.Cmd) {
	cmd := m.dashboard.sidebar.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) DashboardView() string {
	mainScreenStyle := lipgloss.NewStyle()

	sidebarBox := m.dashboard.sidebar.SidebarView(m)

	contentBox := mainScreenStyle.Width(m.widthContainer - 21).Height(m.heightContainer).Render(m.dashboard.sidebar.ViewSelectedTabContent())

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox, contentBox)
}

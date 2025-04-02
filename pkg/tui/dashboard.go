package tui

import (
	tabs "p1/pkg/tabs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	sidebar *Sidebar
}

func NewDashboard() *Dashboard {
	projectsTab := tabs.NewProjectsTab()
	serversTab := tabs.NewServersTab()
	brokersTab := tabs.NewBrokersTab()
	// serversTab := NewTab("servers", NewServerCollection())
	// brokersTab := NewTab("brokers", NewBrokerCollection())
	aboutTab := tabs.NewAboutTab()
	exitTab := tabs.NewExitTab()

	return &Dashboard{
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			aboutTab,
			exitTab,
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
	sidebarBox := m.dashboard.sidebar.SidebarView(m)
	paddedStyle := lipgloss.NewStyle().Padding(1).Border(lipgloss.NormalBorder())

	contentBox := paddedStyle.Width(m.widthContainer - 31).Height(m.heightContainer).Render(m.dashboard.sidebar.ViewSelectedTabContent())

	return lipgloss.JoinHorizontal(lipgloss.Left, sidebarBox, contentBox)
}

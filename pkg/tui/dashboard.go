package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	sidebar    *Sidebar
	mainScreen *MainScreen
	pjc        *ProjectCollection
	sc         *ServerCollection
	bc         *BrokerCollection
}

type MainScreen struct {
	focused bool
	Content Content
}

func NewDashboard() Dashboard {
	projectCollection := NewProjectCollection()
	serverCollection := NewServerCollection()
	brokerCollection := NewBrokerCollection()
	projectsTab := NewTab(projectsScreen, "Projects", false, projectCollection)
	serversTab := NewTab(serversScreen, "Servers", false, serverCollection)
	brokersTab := NewTab(brokersScreen, "Brokers", false, brokerCollection)

	return Dashboard{
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			NewTab(aboutScreen, "About", false, NewAbout())),
		mainScreen: &MainScreen{
			focused: false,
			Content: NewEmptyContent(),
		},
		pjc: projectCollection,
	}
}

func (m model) DashboardSwitch() (model, tea.Cmd) {
	m.dashboard = NewDashboard()
	m = m.SwitchPage(dashboardPage)
	return m, nil
}

func (m model) DashboardUpdate(msg tea.Msg) (model, tea.Cmd) {
	parentMsg := msg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j":
			// Move down in sidebar
			if m.dashboard.sidebar.activeTab < len(m.dashboard.sidebar.tabs)-1 && m.dashboard.sidebar.focused && !m.dashboard.mainScreen.focused {
				m.dashboard.sidebar.activeTab++
			}
		case "k":
			// Move up in sidebar
			if m.dashboard.sidebar.activeTab > 0 && m.dashboard.sidebar.focused && !m.dashboard.mainScreen.focused {
				m.dashboard.sidebar.activeTab--
			}

		case "enter":
			// Focus on main screen
			m.dashboard.mainScreen.focused = true
			m.dashboard.sidebar.focused = false
		case "ctrl+k":
			// Return to sidebar
			m.dashboard.mainScreen.focused = false
			m.dashboard.sidebar.focused = true
		}
	}
	cmd := tea.Batch(m.dashboard.mainScreen.Content.Update(parentMsg), m.dashboard.sidebar.Update(parentMsg))

	return m, cmd
}

func (m model) DashboardView() string {
	mainScreenStyle := lipgloss.NewStyle()

	sidebarBox := m.dashboard.sidebar.SidebarView(m)

	contentBox := mainScreenStyle.Width(m.widthContainer - 21).Height(m.heightContainer).Render(
		m.dashboard.mainScreen.Content.View(),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox, contentBox)
}

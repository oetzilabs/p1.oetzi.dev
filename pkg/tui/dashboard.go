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
	projectsTab := NewTab("projects", "Projects", projectCollection)
	serversTab := NewTab("servers", "Servers", serverCollection)
	brokersTab := NewTab("brokers", "Brokers", brokerCollection)

	return Dashboard{
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			NewTab("about", "About", NewAbout())),
		mainScreen: &MainScreen{
			focused: false,
			Content: NewEmptyContent(),
		},
		pjc: projectCollection,
		sc:  serverCollection,
		bc:  brokerCollection,
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
		case "ctrl+shift+k":
			// Return to sidebar
			m.dashboard.mainScreen.focused = false
			m.dashboard.sidebar.focused = true
		}
	}
	cmd := m.dashboard.mainScreen.Content.Update(parentMsg)
	cmd2 := m.dashboard.sidebar.Update(parentMsg)
	cmd = tea.Batch(cmd, cmd2)

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

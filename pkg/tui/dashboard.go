package tui

import (
	tabs "p1/pkg/tabs"
	"p1/pkg/tui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	sidebar *Sidebar
	theme   *theme.Theme
	width   int
	height  int
}

func NewDashboard(theme *theme.Theme) *Dashboard {
	// main
	projectsTab := tabs.NewProjectsTab()
	serversTab := tabs.NewServersTab()
	brokersTab := tabs.NewBrokersTab()

	// bottom
	aboutTab := tabs.NewAboutTab()
	exitTab := tabs.NewExitTab()

	return &Dashboard{
		theme: theme,
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			aboutTab,
			exitTab,
		),
	}
}

func (d *Dashboard) UpdateSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *Dashboard) Update(msg tea.Msg) tea.Cmd {
	cmd := d.sidebar.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
	}

	return cmd
}

func (d *Dashboard) View() string {
	sidebarBox := d.sidebar.SidebarView()
	paddedStyle := lipgloss.NewStyle().Padding(1)

	contentBox := paddedStyle.Width(d.width - d.sidebar.width).Height(d.height).Render(d.sidebar.ViewSelectedTabContent())

	return lipgloss.JoinHorizontal(lipgloss.Left, sidebarBox, contentBox)
}

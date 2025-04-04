package tui

import (
	"p1/pkg/models"
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
	footer  *models.Footer
}

func NewDashboard(theme *theme.Theme) *Dashboard {
	// main
	projectsTab := tabs.NewProjectsTab()
	serversTab := tabs.NewServersTab()
	brokersTab := tabs.NewBrokersTab()

	// bottom
	aboutTab := tabs.NewAboutTab()
	exitTab := tabs.NewExitTab()

	footerCommands := []models.FooterCommand{
		{Key: "q", Value: "quit"},
		{Key: "←/→", Value: "switch tabs"},
	}

	return &Dashboard{
		theme: theme,
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			aboutTab,
			exitTab,
		),
		footer: models.NewFooter(theme, footerCommands),
	}
}

func (d *Dashboard) UpdateSize(width, height int) {
	d.width = width
	d.height = height
	d.footer.UpdateWidth(width)
}

func (d *Dashboard) Update(msg tea.Msg) tea.Cmd {
	cmd := d.sidebar.Update(msg)
	cmd = tea.Batch(cmd, d.footer.Update(msg))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.footer.UpdateWidth(msg.Width - d.sidebar.width)

	}

	return cmd
}

func (d *Dashboard) View() string {
	sidebarBox := d.sidebar.SidebarView()
	paddedStyle := lipgloss.NewStyle().Padding(1)

	// Subtract some height for the footer
	contentHeight := d.height - lipgloss.Height(d.footer.View())
	contentBox := paddedStyle.Width(d.width - d.sidebar.width).Height(contentHeight).Render(d.sidebar.ViewSelectedTabContent())

	mainContent := lipgloss.JoinVertical(
		lipgloss.Top,
		contentBox,
		d.footer.View(),
	)

	mainContent = lipgloss.JoinHorizontal(lipgloss.Left, sidebarBox, mainContent)

	return mainContent
}

package tui

import (
	"p1/pkg/models"
	tabs "p1/pkg/tabs"
	"p1/pkg/tui/theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	sidebar   *Sidebar
	theme     *theme.Theme
	width     int
	height    int
	viewport  viewport.Model
	footer    *models.Footer
	hasScroll bool
}

var modifiedKeyMap = viewport.KeyMap{
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "½ page down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
}

func NewDashboard(theme *theme.Theme) *Dashboard {
	// main
	projectsTab := tabs.NewProjectsTab()
	serversTab := tabs.NewServersTab()
	brokersTab := tabs.NewBrokersTab()

	// bottom
	aboutTab := tabs.NewAboutTab()
	exitTab := tabs.NewExitTab()

	footerCommands := []models.FooterCommand{}

	footer := models.NewFooter(theme, footerCommands)
	footer.ResetCommands()

	vp := viewport.New(0, 0)
	vp.HighPerformanceRendering = false

	return &Dashboard{
		theme: theme,
		sidebar: NewSidebar(
			projectsTab,
			serversTab,
			brokersTab,
			aboutTab,
			exitTab,
		),
		footer:   footer,
		viewport: vp,
	}
}

func (d *Dashboard) Update(msg tea.Msg) tea.Cmd {
	cmd := d.sidebar.Update(msg)
	cmd = tea.Batch(cmd, d.footer.Update(msg))

	// if d.hasScroll {
	// 	d.width = d.width - 4
	// } else {
	// 	d.width = d.width - 0
	// }

	d.viewport.KeyMap = modifiedKeyMap

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+c":
		// 	return tea.Quit
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.viewport.Width = msg.Width - d.sidebar.width
		d.viewport.Height = msg.Height - lipgloss.Height(d.footer.View())
		d.footer.UpdateWidth(msg.Width - d.sidebar.width)
	}

	d.viewport.SetContent(d.sidebar.ViewSelectedTabContent())
	d.viewport.Width = d.width
	d.viewport.Height = d.height - lipgloss.Height(d.footer.View())
	var cmd2 tea.Cmd
	d.viewport, cmd2 = d.viewport.Update(msg)
	cmd = tea.Batch(cmd, cmd2)
	d.hasScroll = d.viewport.VisibleLineCount() < d.viewport.TotalLineCount()

	return cmd
}

func (d *Dashboard) View() string {
	sidebarBox := d.sidebar.View()

	footerContent := d.footer.View()

	content := d.viewport.View()

	var view string
	if d.hasScroll {
		view = lipgloss.JoinHorizontal(
			lipgloss.Left,
			content,
			d.theme.Base().Width(1).Render(),
			d.getScrollbar(content),
		)
	} else {
		view = lipgloss.JoinHorizontal(
			lipgloss.Left,
			content,
		)
	}

	mainContent := lipgloss.JoinVertical(
		lipgloss.Top,
		view,
		footerContent,
	)

	return lipgloss.JoinHorizontal(lipgloss.Left, sidebarBox, mainContent)
}

func (d *Dashboard) getScrollbar(content string) string {
	y := d.viewport.YOffset
	vh := d.viewport.Height
	ch := lipgloss.Height(content)
	if vh >= ch {
		return ""
	}

	height := (vh * vh) / ch
	maxScroll := ch - vh
	nYP := 1.0 - (float64(y) / float64(maxScroll))
	if nYP <= 0 {
		nYP = 1
	} else if nYP >= 1 {
		nYP = 0
	}

	bar := d.theme.Base().
		Height(height).
		Width(1).
		Background(d.theme.Accent()).
		Render()

	style := d.theme.Base().Width(1).Height(vh)

	return style.Render(lipgloss.PlaceVertical(vh, lipgloss.Position(nYP), bar))
}

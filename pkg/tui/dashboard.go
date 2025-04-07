package tui

import (
	"p1/pkg/client"
	"p1/pkg/interfaces"
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
	client    *client.Client
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
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	),
}

func NewDashboard(theme *theme.Theme, client *client.Client) *Dashboard {
	// main
	projectsTab := tabs.NewProjectsTab(client)
	serversTab := tabs.NewServersTab(client)
	brokersTab := tabs.NewBrokersTab(client)

	// bottom
	aboutTab := tabs.NewAboutTab()
	exitTab := tabs.NewExitTab()

	footerCommands := []interfaces.FooterCommand{}

	footer := models.NewFooter(theme, footerCommands)
	footer.ResetCommands()

	vp := viewport.New(0, 0)
	vp.HighPerformanceRendering = false

	return &Dashboard{
		client: client,
		theme:  theme,
		sidebar: NewSidebar(
			client,
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

	mainScreenContent := d.sidebar.SelectedTabContentView()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+c":
		// 	return tea.Quit
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.updateViewport(d.width-d.sidebar.width, d.height-lipgloss.Height(d.footer.View()))
		d.footer.UpdateWidth(msg.Width - d.sidebar.width)
	}
	d.viewport.SetContent(mainScreenContent)

	var cmd2 tea.Cmd
	d.viewport, cmd2 = d.viewport.Update(msg)
	cmd = tea.Batch(cmd, cmd2)

	d.hasScroll = d.viewport.VisibleLineCount() < d.viewport.TotalLineCount()

	return cmd
}

func (d *Dashboard) updateViewport(width, height int) {
	d.viewport.Width = width
	d.viewport.Height = height
	d.viewport.Style = d.viewport.Style.
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		MaxHeight(height).
		MaxWidth(width)
}

func (d *Dashboard) View() string {
	sidebarBox := d.sidebar.View()

	footerContent := d.footer.View()

	content := d.viewport.View()

	var view string
	view = lipgloss.JoinHorizontal(
		lipgloss.Left,
		content,
		d.getScrollbar(content),
	)
	// if d.hasScroll {
	// } else {
	// 	view = lipgloss.JoinHorizontal(
	// 		lipgloss.Left,
	// 		content,
	// 	)
	// }

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

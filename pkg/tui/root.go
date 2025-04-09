package tui

import (
	"context"
	"strings"

	"p1/pkg/client"
	"p1/pkg/interfaces"
	"p1/pkg/models"
	tabs "p1/pkg/tabs"
	"p1/pkg/tui/theme"
	"p1/pkg/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

type model struct {
	ready          bool
	renderer       *lipgloss.Renderer
	context        context.Context
	theme          theme.Theme
	width          int
	height         int
	client         *client.Client
	viewport       viewport.Model
	footer         *models.Footer
	hasScroll      bool
	tabs           []tabs.Tab
	activeTab      int
	focusedSidebar bool
	search         textinput.Model
	tabsToRender   []tabs.Tab
	sidebarWidth   int
}

func NewModel(renderer *lipgloss.Renderer, client *client.Client) (tea.Model, error) {
	basicTheme := theme.BasicTheme(renderer, nil)

	// main
	projectsTab := tabs.NewProjectsTab()
	serversTab := tabs.NewServersTab()
	brokersTab := tabs.NewBrokersTab()

	// bottom
	aboutTab := tabs.NewAboutTab()
	exitTab := tabs.NewExitTab()

	footerCommands := []interfaces.FooterCommand{}

	footer := models.NewFooter(&basicTheme, footerCommands)
	footer.ResetCommands()

	ti := textinput.New()
	ti.PromptStyle.MaxWidth(23)
	ti.PromptStyle.Width(23)
	ti.Prompt = "# "
	ti.Placeholder = "Search" + strings.Repeat(" ", 23-lipgloss.Width("Search"))

	tabs := []tabs.Tab{
		projectsTab,
		serversTab,
		brokersTab,
		aboutTab,
		exitTab,
	}

	result := model{
		ready:          false,
		renderer:       renderer,
		context:        context.Background(),
		theme:          basicTheme,
		width:          0,
		height:         0,
		client:         client,
		footer:         footer,
		tabs:           tabs,
		activeTab:      0,
		focusedSidebar: true,
		search:         ti,
		tabsToRender:   tabs,
		sidebarWidth:   30,
	}

	return result, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width-m.sidebarWidth, msg.Height-lipgloss.Height(m.footer.View()))
			m.viewport.Style = m.viewport.Style.
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				MaxHeight(m.width - m.sidebarWidth).
				MaxWidth(m.height - lipgloss.Height(m.footer.View()))
		} else {
			m.viewport.Width = m.width - m.sidebarWidth
			m.viewport.Height = m.height - lipgloss.Height(m.footer.View())
		}
		m.footer.UpdateWidth(msg.Width - m.sidebarWidth)

	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			m.search.Focus()
		case "esc":
			if m.search.Focused() {
				m.search.Blur()
			} else {
				m.focusedSidebar = false
			}
		case "ctrl+c":
			if m.search.Focused() {
				m.search.Blur()
			} else {
				return m, tea.Quit
			}
		case "q":
			if !m.search.Focused() {
				return m, tea.Quit
			}
		case "j", "down":
			// Move down in sidebar
			if !m.search.Focused() && m.activeTab < len(m.tabsToRender)-1 && m.focusedSidebar {
				m.activeTab++
			}
		case "k", "up":
			// Move up in sidebar
			if !m.search.Focused() && m.activeTab > 0 && m.focusedSidebar {
				m.activeTab--
			}
		case "enter", "tab":
			// unfocus sidebar
			if m.search.Focused() {
				m.search.Blur()
			}
			m.focusedSidebar = false
		case "ctrl+k":
			// Return to sidebar
			if !m.search.Focused() {
				m.focusedSidebar = true
			}
		}
	}

	cmds = append(cmds, m.footer.Update(msg))

	m.viewport.KeyMap = modifiedKeyMap

	mainScreenContent := m.SelectedTabContentView()

	m.viewport.SetContent(mainScreenContent)

	var cmd2 tea.Cmd
	m.viewport, cmd2 = m.viewport.Update(msg)
	if cmd2 != nil {
		cmds = append(cmds, cmd2)
	}

	tiM, cmd := m.search.Update(msg)
	m.search = tiM
	cmds = append(cmds, cmd)

	m.tabsToRender = utils.FilterTabs(
		m.tabs,
		m.search.Value(),
	)

	if len(m.tabsToRender) == 0 {
		m.activeTab = 0
	}

	tab := m.tabsToRender[m.activeTab]
	tabCmd := tab.Update(msg)
	cmds = append(cmds, tabCmd)

	m.hasScroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string = ""
	paddedStyle := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingTop(1).AlignVertical(lipgloss.Top)
	var sidebar string

	sidebar += m.search.View() + "\n\n"

	mainTabs := utils.FilterTabsGroup(m.tabsToRender, tabs.AlignTop)
	bottomTabs := utils.FilterTabsGroup(m.tabsToRender, tabs.AlignBottom)

	for _, tab := range mainTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()
		isActive := m.tabsToRender[m.activeTab].ID == tab.ID
		sidebar += m.formatTabEntry(display, isActive) + "\n"
	}

	sidebar += "__FILLER_VERTICAL__\n"

	for _, tab := range bottomTabs {
		if tab.Hidden {
			continue
		}
		display := tab.Display()
		isActive := m.tabsToRender[m.activeTab].ID == tab.ID
		sidebar += m.formatTabEntry(display, isActive) + "\n"
	}

	content = paddedStyle.Render(lipgloss.JoinVertical(lipgloss.Top, sidebar))

	content = lipgloss.JoinHorizontal(lipgloss.Top, content)

	// Vertical spacing
	verticalFillerHeight := max(0, m.height-lipgloss.Height(content))
	verticalSpace := strings.Repeat("\n", verticalFillerHeight)
	finalContent := strings.ReplaceAll(content, "__FILLER_VERTICAL__", verticalSpace)
	sidebar = lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Dark: "#111111", Light: "#EEEEEE"}).
		Width(m.sidebarWidth).
		Height(m.height).
		Render(finalContent)

	footerContent := m.footer.View()

	content = m.viewport.View()

	view := lipgloss.JoinHorizontal(
		lipgloss.Left,
		content,
		m.getScrollbar(content),
	)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Top,
		view,
		footerContent,
	)

	return m.renderer.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.theme.Base().
			MaxWidth(m.width).
			MaxHeight(m.height).
			Render(lipgloss.JoinHorizontal(lipgloss.Left, sidebar, mainContent)),
	)
}

func (m *model) getScrollbar(content string) string {
	y := m.viewport.YOffset
	vh := m.viewport.Height
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

	bar := m.theme.Base().
		Height(height).
		Width(1).
		Background(m.theme.Accent()).
		Render()

	style := m.theme.Base().Width(1).Height(vh)

	return style.Render(lipgloss.PlaceVertical(vh, lipgloss.Position(nYP), bar))
}

func (m *model) formatTabEntry(display string, active bool) string {
	parts := strings.SplitN(display, " ", 2)
	title := parts[0]
	info := ""
	if len(parts) > 1 {
		info = parts[1]
	}

	sidebarWidth := m.sidebarWidth - 7 // Account for padding and cursor
	spaceWidth := max(1, sidebarWidth-lipgloss.Width(title)-lipgloss.Width(info))
	spacing := strings.Repeat(" ", spaceWidth)

	content := title + spacing + info
	if active {
		if m.focusedSidebar {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				PaddingRight(2).
				Render("▶ " + content)
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Bold(true).
			PaddingRight(2).
			Render("▶ " + content)
	}
	return lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		Render(content)
}

func (m *model) SelectedTabContentView() string {
	if len(m.tabsToRender) == 0 {
		return "No tabs found"
	}

	content := m.tabsToRender[m.activeTab].Content
	if content == nil {
		return "No content found"
	}

	return content.View()
}

func (m *model) SelectedTabContentCommands() []interfaces.FooterCommand {
	if len(m.tabsToRender) == 0 {
		return []interfaces.FooterCommand{}
	}

	return m.tabsToRender[m.activeTab].Commands()
}

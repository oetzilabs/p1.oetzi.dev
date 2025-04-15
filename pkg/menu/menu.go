package menu

import (
	"p1/pkg/interfaces"
	"p1/pkg/models"
	"p1/pkg/screens"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	items             []*MenuItem
	selectedItemIndex int
	selectedItem      *MenuItem
	focused           bool
	search            textinput.Model
	searchFocused     bool
	width             int
	height            int
}

type MenuItem struct {
	id     string
	title  string
	screen *screens.Screen
}

func NewMenu() *Menu {
	ti := textinput.New()
	ti.PromptStyle.MaxWidth(23)
	ti.PromptStyle.Width(23)
	ti.Prompt = "▶ "
	ti.Placeholder = "Search" + strings.Repeat(" ", 23-lipgloss.Width("Search"))

	return &Menu{
		items:             []*MenuItem{},
		selectedItemIndex: -1,
		selectedItem:      nil,
		focused:           true,
		searchFocused:     false,
		search:            ti,
		width:             30,
		height:            0,
	}
}

func (m *Menu) AddItem(item *MenuItem) *Menu {
	m.items = append(m.items, item)
	m.selectedItemIndex = len(m.items) - 1
	m.selectedItem = item
	return m
}

func (m *Menu) RemoveItem(item *MenuItem) *Menu {
	for i, v := range m.items {
		if v.id == item.id {
			m.items = slices.Delete(m.items, i, i+1)
			m.selectedItemIndex = i - 1
			m.selectedItem = m.items[m.selectedItemIndex]
			break
		}
	}
	return m
}

func (m *Menu) GetCommands() []*interfaces.FooterCommand {
	cmds := []*interfaces.FooterCommand{}

	cmds = append(cmds, m.selectedItem.screen.Commands...)

	return cmds
}

func (m *Menu) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}
	parentMsg := msg
	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			if m.focused && !m.search.Focused() {
				m.search.Focus()
			}
		case "esc":
			if m.search.Focused() {
				m.search.Blur()
			} else {
				m.focused = false
			}
		case "ctrl+c":
			if m.search.Focused() {
				m.search.Blur()
			}
		case "j", "down":
			if m.focused && !m.search.Focused() && m.selectedItemIndex < len(m.items)-1 {
				m.selectedItemIndex = max(len(m.items)-1, m.selectedItemIndex+1)
				m.selectedItem = m.items[m.selectedItemIndex]
				m.selectedItem.screen.SetFocused(false)
			}
		case "k", "up":
			if m.focused && !m.search.Focused() && m.selectedItemIndex > 0 {
				m.selectedItemIndex = max(0, m.selectedItemIndex-1)
				m.selectedItem = m.items[m.selectedItemIndex]
				m.selectedItem.screen.SetFocused(false)
			}
		case "enter", "tab":
			// unfocus sidebar
			if m.focused && m.search.Focused() {
				m.search.Blur()
			} else if m.focused && !m.search.Focused() {
				m.focused = false
				m.selectedItem.screen.SetFocused(true)
			}
		case "ctrl+k":
			// Return to sidebar
			if !m.focused && !m.search.Focused() {
				m.focused = true
				m.selectedItem.screen.SetFocused(false)
			}
		case "q":
			if m.focused && !m.search.Focused() {
				cmds = append(cmds, tea.Quit)
			}
		}
	}

	tiM, cmd := m.search.Update(parentMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m.search = tiM

	for _, item := range m.items {
		cmds = append(cmds, item.Update(parentMsg))
	}

	return tea.Batch(cmds...)
}

func (m *Menu) View() string {
	bgColor := lipgloss.Color("#222222")

	menuItemStyle := lipgloss.NewStyle().
		Background(bgColor).Width(m.width)

	menuStyle := lipgloss.NewStyle().
		Background(bgColor).Padding(1).Width(m.width)

	var content string

	// Add the search bar
	content += m.search.View() + "\n\n"

	// Add menu items
	for itemIndex, item := range m.items {
		var newLine string
		if itemIndex < len(m.items)-1 {
			newLine = "\n"
		} else {
			newLine = ""
		}
		content += menuItemStyle.Render(m.formatItem(item.View(), itemIndex == m.selectedItemIndex)) + newLine
	}

	var fillercontent string
	// Calculate remaining height and add filler
	currentHeight := lipgloss.Height(content)
	if currentHeight < m.height {
		fillerHeight := m.height - currentHeight - 4
		filler := strings.Repeat(strings.Repeat(" ", m.width)+"\n", fillerHeight)
		fillercontent += menuItemStyle.Render(filler) + "\n"
	}

	return menuStyle.Render(lipgloss.JoinVertical(lipgloss.Top, content, fillercontent))
}

func (m *Menu) Screen() string {
	if m.selectedItem == nil {
		return "Please select a menu item\n" + strings.Repeat("\n", max(0, m.height-2))
	}

	return m.selectedItem.screen.View()
}

func NewMenuItem(id string, title string, screen *screens.Screen) *MenuItem {
	return &MenuItem{
		id:     id,
		title:  title,
		screen: screen,
	}
}

func (m *Menu) GetWidth() int {
	return lipgloss.Width(m.View())
}

func (m *Menu) IsFocused() bool {
	return m.focused
}

func (m *Menu) formatItem(display string, active bool) string {
	parts := strings.SplitN(display, " ", 2)
	title := parts[0]
	sidebarWidth := m.width
	spaceWidth := max(1, sidebarWidth-lipgloss.Width(title))
	spacing := strings.Repeat(" ", spaceWidth)

	content := title + spacing
	if active {
		if m.focused {
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

func (mi *MenuItem) Update(msg tea.Msg) tea.Cmd {
	if mi.screen != nil {
		cmd := mi.screen.Update(msg)
		return cmd
	}
	return nil
}

func (mi *MenuItem) View() string {
	content := mi.title
	return content
}

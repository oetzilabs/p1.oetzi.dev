package menu

import (
	"log/slog"
	"p1/pkg/models"
	"p1/pkg/screens"
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
		selectedItemIndex: 0,
		focused:           true,
		searchFocused:     false,
		search:            ti,
		width:             30,
		height:            0,
	}
}

func (m *Menu) AddItem(item *MenuItem) *Menu {
	m.items = append(m.items, item)
	return m
}

func (m *Menu) RemoveItem(item *MenuItem) *Menu {
	m.items = append(m.items, item)
	return m
}

func (m *Menu) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}
	parentMsg := msg
	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
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
			if m.focused && m.selectedItemIndex < len(m.items)-1 {
				m.selectedItemIndex++
			}
		case "k", "up":
			if m.focused && m.selectedItemIndex > 0 {
				m.selectedItemIndex--
			}
		case "enter", "tab":
			// unfocus sidebar
			if m.focused && m.search.Focused() {
				m.search.Blur()
			}
			m.focused = false
			m.selectedItem = m.items[m.selectedItemIndex]
		case "ctrl+k":
			// Return to sidebar
			if !m.focused && !m.search.Focused() {
				m.focused = true
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

	if m.selectedItem != nil {
		m.selectedItem.screen.SetConstrains(m.width, m.height)
		cmds = append(cmds, m.selectedItem.screen.Update(parentMsg))
	}

	return tea.Batch(cmds...)
}

func (m *Menu) View() string {
	var content string

	// Add the search bar
	content += m.search.View() + "\n"

	// Add menu items
	for itemIndex, item := range m.items {
		content += m.formatItem(item.title+"\n", itemIndex == m.selectedItemIndex)
	}

	// Calculate remaining height and add filler
	currentHeight := lipgloss.Height(content)
	if currentHeight < m.height {
		slog.Info("adding filler", "height", m.height, "currentHeight", currentHeight)
		fillerHeight := m.height - currentHeight
		filler := strings.Repeat(" \n", fillerHeight)
		content += filler
	}

	return lipgloss.JoinVertical(lipgloss.Top, content)
}

func (m *Menu) Screen() string {
	if m.selectedItem == nil {
		return ""
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
	return m.width
}

func (m *Menu) IsFocused() bool {
	return m.focused
}

func (m *Menu) formatItem(display string, active bool) string {
	parts := strings.SplitN(display, " ", 2)
	title := parts[0]
	info := ""
	if len(parts) > 1 {
		info = parts[1]
	}

	sidebarWidth := m.width - 7 // Account for padding and cursor
	spaceWidth := max(1, sidebarWidth-lipgloss.Width(title)-lipgloss.Width(info))
	spacing := strings.Repeat(" ", spaceWidth)

	content := title + spacing + info
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

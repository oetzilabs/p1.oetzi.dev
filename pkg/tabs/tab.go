package tabs

import (
	interfaces "p1/pkg/interfaces"

	tea "github.com/charmbracelet/bubbletea"
)

// Tab represents a tab in the sidebar
type Tab struct {
	ID      string
	Hidden  bool
	Group   TabGroup
	Content interfaces.Content
}

type UpdateTabDisplay struct {
	DisplayLeft  string
	DisplayRight []string
}

// NewTab creates a new tab
func NewTab(id string, content interfaces.Content) Tab {
	return Tab{
		ID:      id,
		Content: content,
	}
}

// Update updates the tab's state
func (t *Tab) Update(msg tea.Msg) tea.Cmd {
	if t.Content == nil {
		return nil
	}

	return t.Content.Update(msg)
}

// View returns the tab's view
func (t *Tab) View() string {
	if t.Content == nil {
		return "No Content Set - This should not happen."
	}

	return t.Content.View()
}

// View returns the tab's view
func (t *Tab) Display() string {
	return t.Content.Display()
}

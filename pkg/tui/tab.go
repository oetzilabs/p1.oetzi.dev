package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Tab represents a tab in the sidebar
type Tab struct {
	ID      string
	Hidden  bool
	Content Content
}

type UpdateTabDisplay struct {
	DisplayLeft  string
	DisplayRight []string
}

// NewTab creates a new tab
func NewTab(id string, content Content) Tab {
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

	cmd := t.Content.Update(msg)

	return cmd
}

// View returns the tab's view
func (t *Tab) View() string {
	return t.Content.View()
}

// View returns the tab's view
func (t *Tab) Display() string {
	return t.Content.Display()
}

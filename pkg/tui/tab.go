package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Tab represents a tab in the sidebar
type Tab struct {
	ID         string
	Hidden     bool
	TabDisplay string
	Content    Content
}

type UpdateTabDisplay struct {
	DisplayLeft  string
	DisplayRight []string
}

// NewTab creates a new tab
func NewTab(id string, name string, content Content) Tab {
	return Tab{
		ID:         id,
		Content:    content,
		TabDisplay: name,
	}
}

// Update updates the tab's state
func (t *Tab) Update(msg tea.Msg) tea.Cmd {
	if t.Content == nil {
		return nil
	}

	cmd := t.Content.Update(msg)
	t.TabDisplay = t.Content.View()

	return cmd
}

// View returns the tab's view
func (t *Tab) View() string {
	return t.TabDisplay
}

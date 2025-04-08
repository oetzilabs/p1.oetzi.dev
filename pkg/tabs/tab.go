package tabs

import (
	interfaces "p1/pkg/interfaces"
	"p1/pkg/messages"
	"p1/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
)

// Tab represents a tab in the sidebar
type Tab struct {
	ID           string
	Hidden       bool
	Group        TabGroup
	Content      interfaces.Content
	Helper       string
	IgnoreSearch bool
}

// NewTab creates a new tab
func NewTab(id string, content interfaces.Content, helper string) Tab {
	return Tab{
		ID:           id,
		Content:      content,
		Helper:       helper,
		IgnoreSearch: false,
	}
}

// Update updates the tab's state
func (t *Tab) Update(msg tea.Msg) tea.Cmd {
	if t.Content == nil {
		return nil
	}
	cmd := t.Content.Update(msg)

	updateHelperMsg := func() tea.Msg { return models.FooterUpdate{Content: t.Helper, Commands: t.Commands()} }
	cmd = tea.Batch(cmd, updateHelperMsg)

	return cmd
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

// Commands returns the tab's commands
func (t *Tab) Commands() []interfaces.FooterCommand {
	if t.Content == nil {
		return []interfaces.FooterCommand{}
	}

	return t.Content.Commands()
}

func (t *Tab) SendMessage(msg messages.Message) {
	if t.Content == nil {
		return
	}
	t.Content.Update(msg)
}

package tabs

import (
	collections "p1/pkg/tui/collections"

	tea "github.com/charmbracelet/bubbletea"
)

type ServersTab struct {
	collection *collections.ServerCollection
}

func NewServersTab() Tab {

	return Tab{
		ID:     "Servers",
		Hidden: false,
		Content: &ServersTab{
			collection: collections.NewServerCollection(),
		},
	}
}

func (pt *ServersTab) Update(msg tea.Msg) tea.Cmd {
	return pt.collection.Update(msg)
}

func (pt *ServersTab) View() string {
	return pt.collection.View()
}

func (pt *ServersTab) Display() string {
	return pt.collection.Display()
}

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
		ID:     "servers",
		Hidden: false,
		Group:  AlignTop,
		Content: &ServersTab{
			collection: collections.NewServerCollection(),
		},
		Helper: "Here you can see all the servers that are currently connected to the network.",
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

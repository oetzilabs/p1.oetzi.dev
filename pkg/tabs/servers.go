package tabs

import (
	"p1/pkg/client"
	"p1/pkg/interfaces"
	collections "p1/pkg/tui/collections"

	tea "github.com/charmbracelet/bubbletea"
)

type ServersTab struct {
	collection *collections.ServerCollection
}

func NewServersTab(client *client.Client) Tab {
	return Tab{
		ID:     "servers",
		Hidden: false,
		Group:  AlignTop,
		Content: &ServersTab{
			collection: collections.NewServerCollection(client),
		},
		Helper: "Here you can see all the servers that are currently connected to the network.",
	}
}

func (st *ServersTab) Update(msg tea.Msg) tea.Cmd {
	return st.collection.Update(msg)
}

func (st *ServersTab) View() string {
	return st.collection.View()
}

func (st *ServersTab) Display() string {
	count := st.collection.Count()
	return "Servers " + count
}

func (st *ServersTab) Commands() []interfaces.FooterCommand {
	return st.collection.Commands()
}

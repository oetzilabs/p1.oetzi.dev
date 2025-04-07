package tabs

import (
	"p1/pkg/client"
	"p1/pkg/interfaces"
	collections "p1/pkg/tui/collections"

	tea "github.com/charmbracelet/bubbletea"
)

type BrokersTab struct {
	collection *collections.BrokerCollection
}

func NewBrokersTab(client *client.Client) Tab {

	return Tab{
		ID:     "brokers",
		Hidden: false,
		Group:  AlignTop,
		Content: &BrokersTab{
			collection: collections.NewBrokerCollection(client),
		},
		Helper: "Here you can see all the brokers that are currently connected to the network.",
	}
}

func (pt *BrokersTab) Update(msg tea.Msg) tea.Cmd {
	return pt.collection.Update(msg)
}

func (pt *BrokersTab) View() string {
	return pt.collection.View()
}

func (pt *BrokersTab) Display() string {
	count := pt.collection.Count()
	return "Brokers " + count
}

func (pt *BrokersTab) Commands() []interfaces.FooterCommand {
	return pt.collection.Commands()
}

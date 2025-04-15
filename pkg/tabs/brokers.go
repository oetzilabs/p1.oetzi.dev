package tabs

import (
	"p1/pkg/interfaces"
	collections "p1/pkg/tui/collections"

	tea "github.com/charmbracelet/bubbletea"
)

type BrokersTab struct {
	collection *collections.BrokerCollection
}

func NewBrokersTab() Tab {

	return Tab{
		ID:     "brokers",
		Hidden: false,
		Group:  AlignTop,
		Content: &BrokersTab{
			collection: collections.NewBrokerCollection(),
		},
		Helper: "Here you can see all the brokers that are currently connected to the network.",
	}
}

func (bt *BrokersTab) Update(msg tea.Msg) tea.Cmd {
	return bt.collection.Update(msg)
}

func (bt *BrokersTab) View() string {
	return bt.collection.View()
}

func (bt *BrokersTab) Display() string {
	count := bt.collection.Count()
	return "Brokers " + count
}

func (bt *BrokersTab) Commands() []*interfaces.FooterCommand {
	return bt.collection.Commands()
}

package tabs

import (
	collections "p1/pkg/tui/collections"

	tea "github.com/charmbracelet/bubbletea"
)

type BrokersTab struct {
	collection *collections.BrokerCollection
}

func NewBrokersTab() Tab {

	return Tab{
		ID:     "Brokers",
		Hidden: false,
		Content: &BrokersTab{
			collection: collections.NewBrokerCollection(),
		},
	}
}

func (pt *BrokersTab) Update(msg tea.Msg) tea.Cmd {
	return pt.collection.Update(msg)
}

func (pt *BrokersTab) View() string {
	return pt.collection.View()
}

func (pt *BrokersTab) Display() string {
	return pt.collection.Display()
}

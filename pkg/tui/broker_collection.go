package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type BrokerCollection struct {
	brokers  []*BrokerView
	selected int
}

func NewBrokerCollection() *BrokerCollection {
	return &BrokerCollection{
		brokers: []*BrokerView{},
	}
}

func (bc *BrokerCollection) AddBroker(broker *BrokerView) {
	bc.brokers = append(bc.brokers, broker)
}

func (bc *BrokerCollection) SelectBroker(id string) {
	for i, broker := range bc.brokers {
		if broker.Id == id {
			bc.selected = i
			return
		}
	}
}

func (bc *BrokerCollection) RemoveBroker(id string) {
	for i, broker := range bc.brokers {
		if broker.Id == id {
			bc.brokers = append(bc.brokers[:i], bc.brokers[i+1:]...)
			return
		}
	}
}

func (bc *BrokerCollection) Update(msg tea.Msg) tea.Cmd {
	parentMsg := msg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		}
	}
	var cmd tea.Cmd

	if len(bc.brokers) > 0 {
		cmd = tea.Batch(bc.brokers[bc.selected].Update(parentMsg), cmd)
	}

	return cmd
}

func (bc *BrokerCollection) View() string {
	if len(bc.brokers) == 0 {
		return "No brokers available. Press 'n' to add a new broker."
	}
	content := bc.brokers[bc.selected].View()

	return content
}

func (bc *BrokerCollection) Display() string {
	if len(bc.brokers) == 0 {
		return "Brokers"
	}
	return fmt.Sprintf("Brokers (%d)", len(bc.brokers))
}

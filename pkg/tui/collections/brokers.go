package tui

import (
	"fmt"
	"p1/pkg/interfaces"
	"p1/pkg/messages"
	models "p1/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
)

type BrokerCollection struct {
	brokers  []*models.Broker
	selected int
}

func NewBrokerCollection() *BrokerCollection {
	return &BrokerCollection{
		brokers: []*models.Broker{},
	}
}

func (bc *BrokerCollection) AddBroker(broker *models.Broker) {
	bc.brokers = append(bc.brokers, broker)
}

func (bc *BrokerCollection) SelectBroker(id string) {
	for i, broker := range bc.brokers {
		if broker.ID == id {
			bc.selected = i
			return
		}
	}
}

func (bc *BrokerCollection) RemoveBroker(id string) {
	for i, broker := range bc.brokers {
		if broker.ID == id {
			bc.brokers = append(bc.brokers[:i], bc.brokers[i+1:]...)
			return
		}
	}
}

func (bc *BrokerCollection) Update(msg tea.Msg) tea.Cmd {
	parentMsg := msg
	switch msg := msg.(type) {
	case messages.Message:
		switch msg.Type {
		case messages.TypeListBrokers:
			bc.brokers = msg.Payload.([]*models.Broker)
		case messages.TypeRegisterBroker:
			bc.AddBroker(msg.Payload.(*models.Broker))
		case messages.TypeRemoveBroker:
			bc.RemoveBroker(msg.Payload.(string))
		}
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

func (bc *BrokerCollection) Count() string {
	return fmt.Sprintf("(%d)", len(bc.brokers))
}

func (bc *BrokerCollection) Commands() []*interfaces.FooterCommand {
	return []*interfaces.FooterCommand{
		{Key: "ctrl+n", Value: "New Broker"},
	}
}

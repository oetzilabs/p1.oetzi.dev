package tui

import tea "github.com/charmbracelet/bubbletea"

type BrokerCollection struct {
	brokers        []*BrokerView
	selected       int
	to_remove      string
	confirm_delete *Confirmation
	dialog         *Dialog
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

func (bc *BrokerCollection) ConfirmRemoveBroker(id string) {
	bc.to_remove = id
	name := bc.brokers[bc.selected].Name
	confirm := NewConfirmation("Do you really wish to delete "+name+"?", func() {
		bc.RemoveBroker(bc.to_remove)
	}, func() {
		bc.to_remove = ""
	})
	bc.confirm_delete = confirm
}

func (bc *BrokerCollection) AddBrokerDialog() {
	inputs := []Input{*NewInput("ID", "123"), *NewInput("Name", "My Broker"), *NewInput("URL", "mqtt://example.com")}
	dialog := NewDialog("Enter the broker details", func(values ...string) {
		bc.AddBroker(NewBroker(values[0], values[1], values[2]))
	}, func() {
		bc.dialog = nil
	}, inputs...)
	bc.dialog = dialog
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
	if bc.confirm_delete != nil {
		cmd = bc.confirm_delete.Update(msg)
	}

	if len(bc.brokers) > 0 {
		cmd = tea.Batch(bc.brokers[bc.selected].Update(parentMsg), cmd)
	}
	return cmd
}

func (bc *BrokerCollection) View() string {
	content := bc.brokers[bc.selected].View()
	if bc.to_remove != "" {
		content += bc.confirm_delete.View()
	}
	return content
}

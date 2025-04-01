package tui

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

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
	inputs := []Input{
		*NewInput(
			"Name",
			"",
			&InputOptions{
				Focused:     true,
				Disabled:    false,
				Placeholder: "Broker Name",
			},
		),
		*NewInput(
			"URL",
			"",
			&InputOptions{
				Focused:     false,
				Disabled:    false,
				Placeholder: "mqtt://example.com",
			},
		),
	}
	dialog := NewDialog(
		"Enter the broker details",
		inputs,
		func(values interface{}) {
			v := reflect.ValueOf(values)
			name := v.FieldByName("Name").String()
			url := v.FieldByName("URL").String()
			bc.AddBroker(NewBroker(uuid.NewString(), name, url))
			bc.dialog = nil
		},
		func() {
			bc.dialog = nil
		},
	)
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

	// Send display information up to parent
	updateCmd := tea.Cmd(func() tea.Msg {
		return UpdateTabDisplay{
			DisplayLeft:  "Brokers",
			DisplayRight: []string{fmt.Sprintf("(%d)", len(bc.brokers))},
		}
	})
	cmd = tea.Batch(cmd, updateCmd)
	return cmd
}

func (bc *BrokerCollection) View() string {
	if len(bc.brokers) == 0 {
		return "No brokers available. Press 'n' to add a new broker."
	}
	content := bc.brokers[bc.selected].View()
	if bc.to_remove != "" {
		content += bc.confirm_delete.View()
	}
	return content
}

func (bc *BrokerCollection) Display() string {
	if len(bc.brokers) == 0 {
		return "Brokers"
	}
	return fmt.Sprintf("Brokers (%d)", len(bc.brokers))
}

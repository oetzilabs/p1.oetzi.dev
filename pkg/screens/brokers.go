package screens

import (
	"fmt"
	"p1/pkg/models"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BrokersScreen struct {
	collection []*models.Broker
	selected   int
}

func NewBrokersScreen(renderer *lipgloss.Renderer) *Screen {
	screen := &BrokersScreen{
		collection: []*models.Broker{
			models.NewBroker("Broker 1", "http://localhost:8080"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
			models.NewBroker("Broker 2", "http://localhost:8081"),
		},
		selected: 0,
	}
	return NewScreen(renderer, screen)
}

func (bs *BrokersScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for _, broker := range bs.collection {
		cmds = append(cmds, broker.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (bs *BrokersScreen) View() string {
	content := fmt.Sprintf("Brokers (%d)\n", len(bs.collection))
	for _, broker := range bs.collection {
		content += broker.View() + "\n"
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, content)
}

func (s *BrokersScreen) AddBroker(broker *models.Broker) *BrokersScreen {
	s.collection = append(s.collection, broker)
	return s
}

func (s *BrokersScreen) RemoveBroker(broker *models.Broker) *BrokersScreen {
	for i, p := range s.collection {
		if p.ID == broker.ID {
			s.collection = slices.Delete(s.collection, i, i+1)
			break
		}
	}
	return s
}

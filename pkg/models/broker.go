package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Broker struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewBroker(name string, url string) Broker {
	id := uuid.New().String()
	return Broker{
		ID:   id,
		Name: name,
		URL:  url,
	}
}

func (p *Broker) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Broker
	return nil
}

func (p *Broker) View() string {
	mainStyle := lipgloss.NewStyle().Padding(2)
	return mainStyle.Render(p.Name)
}

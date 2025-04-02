package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Actor struct {
	ID string `json:"id"`
}

func NewActor() Actor {
	id := uuid.New().String()
	return Actor{
		ID: id,
	}
}

func (p *Actor) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Actor
	return nil
}

func (p *Actor) View() string {
	mainStyle := lipgloss.NewStyle().Padding(2)
	return mainStyle.Render(p.ID)
}

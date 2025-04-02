package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Server struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewServer(name string, url string) Server {
	id := uuid.New().String()
	return Server{
		ID:   id,
		Name: name,
		URL:  url,
	}
}

func (p *Server) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Server
	return nil
}

func (p *Server) View() string {
	mainStyle := lipgloss.NewStyle().Padding(2)
	return mainStyle.Render(p.Name)
}

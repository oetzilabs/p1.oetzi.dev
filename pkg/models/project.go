package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewProject(id string, name string) *Project {
	return &Project{
		ID:   id,
		Name: name,
	}
}

func (p *Project) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Project
	return nil
}

func (p *Project) View() string {
	mainStyle := lipgloss.NewStyle()
	return mainStyle.Render(p.Name)
}

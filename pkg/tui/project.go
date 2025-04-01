package tui

import (
	"p1/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectView struct {
	*models.Project
}

func NewProject(id string, name string) *ProjectView {
	return &ProjectView{
		Project: &models.Project{
			ID:   id,
			Name: name,
		},
	}
}

func (p *ProjectView) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Project
	return nil
}

func (p *ProjectView) View() string {
	mainStyle := lipgloss.NewStyle().Padding(2)
	return mainStyle.Render(p.Name)
}

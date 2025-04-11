package screens

import (
	"fmt"
	"p1/pkg/models"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectsScreen struct {
	collection []*models.Project
	selected   int
}

func NewProjectsScreen(renderer *lipgloss.Renderer) *Screen {
	screen := &ProjectsScreen{
		collection: []*models.Project{},
		selected:   0,
	}
	return NewScreen(renderer, screen)
}

func (ps *ProjectsScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for _, project := range ps.collection {
		cmds = append(cmds, project.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (ps *ProjectsScreen) View() string {
	content := fmt.Sprintf("Projects (%d))\n", len(ps.collection))

	for _, project := range ps.collection {
		content += project.View() + "\n"
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, content)
}

func (s *ProjectsScreen) AddProject(project *models.Project) *ProjectsScreen {
	s.collection = append(s.collection, project)
	return s
}

func (s *ProjectsScreen) RemoveProject(project *models.Project) *ProjectsScreen {
	for i, p := range s.collection {
		if p.ID == project.ID {
			s.collection = slices.Delete(s.collection, i, i+1)
			break
		}
	}
	return s
}

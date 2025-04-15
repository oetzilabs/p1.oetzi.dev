package tui

import (
	"fmt"
	"p1/pkg/interfaces"
	"p1/pkg/messages"
	models "p1/pkg/models"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Projects struct {
	projects    []*models.Project
	selected    int
	placeholder string
}

func NewProjectCollection() *Projects {
	return &Projects{
		projects:    []*models.Project{},
		selected:    0,
		placeholder: "There are no projects yet",
	}
}

func (pc *Projects) AddProject(project *models.Project) {
	pc.projects = append(pc.projects, project)
}

func (pc *Projects) SelectProject(id string) {
	for i, project := range pc.projects {
		if project.ID == id {
			pc.selected = i
			return
		}
	}
}

func (pc *Projects) RemoveProject(id string) {
	for i, project := range pc.projects {
		if project.ID == id {
			pc.projects = slices.Delete(pc.projects, i, i+1)
			return
		}
	}
}

func (pc *Projects) Update(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {
	case messages.Message:
		switch msg.Type {
		case messages.TypeListProjects:
			pc.projects = msg.Payload.([]*models.Project)
		case messages.TypeRegisterProjects:
			pc.AddProject(msg.Payload.(*models.Project))
		case messages.TypeRemoveProjects:
			pc.RemoveProject(msg.Payload.(string))
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
		}
	}
	var cmd tea.Cmd

	for _, project := range pc.projects {
		// Update both the tab and its content
		if pcmd := project.Update(msg); pcmd != nil {
			cmd = tea.Batch(cmd, pcmd)
		}

	}

	return cmd
}

func (pc *Projects) View() string {
	if len(pc.projects) == 0 {
		return pc.placeholder
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, pc.projects[pc.selected].View())
}

func (pc *Projects) Count() string {
	return fmt.Sprintf("(%d)", len(pc.projects))
}

func (pc *Projects) Commands() []*interfaces.FooterCommand {
	return []*interfaces.FooterCommand{
		{Key: "ctrl+n", Value: "New Project"},
	}
}

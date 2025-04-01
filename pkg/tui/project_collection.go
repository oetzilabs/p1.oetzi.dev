package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ProjectCollection struct {
	projects []*ProjectView
	selected int
}

func NewProjectCollection() *ProjectCollection {
	return &ProjectCollection{
		projects: []*ProjectView{},
		selected: 0,
	}
}

func (pc *ProjectCollection) AddProject(project *ProjectView) {
	pc.projects = append(pc.projects, project)
}

func (pc *ProjectCollection) SelectProject(id string) {
	for i, project := range pc.projects {
		if project.ID == id {
			pc.selected = i
			return
		}
	}
}

func (pc *ProjectCollection) RemoveProject(id string) {
	for i, project := range pc.projects {
		if project.ID == id {
			pc.projects = append(pc.projects[:i], pc.projects[i+1:]...)
			return
		}
	}
}

func (pc *ProjectCollection) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
		case "d":
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

func (pc *ProjectCollection) View() string {
	if len(pc.projects) == 0 {
		return "No projects available. Press 'n' to add a new project."
	}

	return pc.projects[pc.selected].View()
}

func (pc *ProjectCollection) Display() string {
	if len(pc.projects) == 0 {
		return "Projects"
	}
	return fmt.Sprintf("Projects (%d)", len(pc.projects))
}

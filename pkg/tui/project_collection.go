package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectCollection struct {
	projects       []*Project
	selected       int
	to_remove      string
	confirm_delete *Confirmation
}

func NewProjectCollection() *ProjectCollection {
	return &ProjectCollection{
		projects: []*Project{},
		selected: 0,
	}
}

func (pc *ProjectCollection) AddProject(project *Project) {
	pc.projects = append(pc.projects, project)
}

func (pc *ProjectCollection) SelectProject(id string) {
	for i, project := range pc.projects {
		if project.id == id {
			pc.selected = i
			return
		}
	}
}

func (pc *ProjectCollection) ConfirmRemoveProject(id string) {
	pc.to_remove = id
	name := pc.projects[pc.selected].name
	confirm := NewConfirmation("Do you really wish to delete"+name+"?", func() {
		pc.RemoveProject(pc.to_remove)
	}, func() {
		pc.to_remove = ""
	})
	pc.confirm_delete = confirm
}

func (pc *ProjectCollection) RemoveProject(id string) {
	for i, project := range pc.projects {
		if project.id == id {
			pc.projects = append(pc.projects[:i], pc.projects[i+1:]...)
			return
		}
	}
}

func (pc *ProjectCollection) Update(msg tea.Msg) tea.Cmd {
	parentMsg := msg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// New project
			pid := "123"
			pc.AddProject(NewProject(pid, "Test"))
			pc.SelectProject(pid)
		case "d":
			// confirm project deletion
			pc.ConfirmRemoveProject(pc.projects[pc.selected].id)
		}
	}
	var cmd tea.Cmd
	if pc.confirm_delete != nil {
		cmd = pc.confirm_delete.Update(msg)
	}

	if len(pc.projects) > 0 {
		cmd = tea.Batch(pc.projects[pc.selected].Update(parentMsg), cmd)
	}
	return cmd
}

func (pc *ProjectCollection) View() string {
	content := pc.projects[pc.selected].View()
	if pc.to_remove != "" {
		content += pc.confirm_delete.View()
	}
	return content
}

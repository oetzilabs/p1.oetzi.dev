package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectCollection struct {
	projects       []*ProjectView
	selected       int
	to_remove      string
	confirm_delete *Confirmation
	dialog         *Dialog
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

func (pc *ProjectCollection) ConfirmRemoveProject(id string) {
	pc.to_remove = id
	name := pc.projects[pc.selected].Name
	confirm := NewConfirmation("Do you really wish to delete"+name+"?", func() {
		pc.RemoveProject(pc.to_remove)
	}, func() {
		pc.to_remove = ""
	})
	pc.confirm_delete = confirm
}

func (pc *ProjectCollection) AddProjectDialog() {
	inputs := []Input{*NewInput("ID", "123"), *NewInput("Name", "My Project")}
	dialog := NewDialog("Enter the name of the project", func(values ...string) {
		pc.AddProject(NewProject(values[0], values[1]))
	}, func() {
		pc.dialog = nil
	}, inputs...)
	pc.dialog = dialog
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
	parentMsg := msg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// TODO: New project dialog window
			pc.AddProjectDialog()
		case "d":
			// confirm project deletion
			pc.ConfirmRemoveProject(pc.projects[pc.selected].ID)
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

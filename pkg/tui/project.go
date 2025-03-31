package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Project struct {
	id   string
	name string
}

func NewProject(id string, name string) *Project {
	return &Project{
		id:   id,
		name: name,
	}
}

func (p *Project) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Project
	return nil
}

func (p *Project) View() string {
	var content string = "Project Content"
	return content
}

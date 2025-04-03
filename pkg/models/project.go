package models

import (
	"encoding/json"
	"p1/pkg/api"

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
	switch msg := msg.(type) {
	case api.WebSocketDataUpdate:
		if msg.Type == "project" {
			var data Project
			if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
				return tea.Cmd(func() tea.Msg {
					return VisibleError{
						Message: err.Error(),
					}
				})
			}
			p.ID = data.ID
			p.Name = data.Name
		}
	}
	return nil
}

func (p *Project) View() string {
	mainStyle := lipgloss.NewStyle()
	return mainStyle.Render(p.Name)
}

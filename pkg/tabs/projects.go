package tabs

import (
	collections "p1/pkg/tui/collections"

	tea "github.com/charmbracelet/bubbletea"
)

type ProjectsTab struct {
	collection *collections.Projects
}

func NewProjectsTab() Tab {

	return Tab{
		ID:     "projects",
		Hidden: false,
		Group:  TabGroupsMain,
		Content: &ProjectsTab{
			collection: collections.NewProjectCollection(),
		},
	}
}

func (pt *ProjectsTab) Update(msg tea.Msg) tea.Cmd {
	return pt.collection.Update(msg)
}

func (pt *ProjectsTab) View() string {
	return pt.collection.View()
}

func (pt *ProjectsTab) Display() string {
	return pt.collection.Display()
}

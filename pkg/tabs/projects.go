package tabs

import (
	"p1/pkg/interfaces"
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
		Group:  AlignTop,
		Content: &ProjectsTab{
			collection: collections.NewProjectCollection(),
		},
		Helper: "Here you can see all the projects that are on your network.",
	}
}

func (pt *ProjectsTab) Update(msg tea.Msg) tea.Cmd {
	return pt.collection.Update(msg)
}

func (pt *ProjectsTab) View() string {
	return pt.collection.View()
}

func (pt *ProjectsTab) Display() string {
	count := pt.collection.Count()
	return "Projects " + count
}

func (pt *ProjectsTab) Commands() []interfaces.FooterCommand {
	return pt.collection.Commands()
}

package screens

import (
	"fmt"
	"p1/pkg/dialog"
	"p1/pkg/interfaces"
	"p1/pkg/models"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type ProjectViewStatus string

const (
	ProjectViewStatusList ProjectViewStatus = "list"
	ProjectViewStatusNew  ProjectViewStatus = "new"
	ProjectViewStatusEdit ProjectViewStatus = "edit"
)

type ProjectsScreen struct {
	collection []*models.Project
	selected   int
	dialog     *dialog.Dialog
	pdb        *ProjectDialogBody
	viewstatus ProjectViewStatus
}

var (
	newProjectKey = &interfaces.FooterCommand{Key: "ctrl+n", Value: "New Project"}
)

func NewProjectsScreen(renderer *lipgloss.Renderer) *Screen {
	pdb := NewProjectDialogBody()
	screen := &ProjectsScreen{
		collection: []*models.Project{},
		selected:   0,
		dialog:     dialog.NewDialog("What is the name of your project?", pdb),
		pdb:        pdb,
		viewstatus: ProjectViewStatusList,
	}
	return NewScreen(renderer, screen, newProjectKey)
}

func (ps *ProjectsScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for _, project := range ps.collection {
		cmds = append(cmds, project.Update(msg))
	}
	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		ps.dialog.UpdateSize(msg.Width-msg.MenuWidth-8, msg.Height-msg.FooterHeight)
	case tea.KeyMsg:
		switch msg.String() {
		case newProjectKey.Key:
			ps.viewstatus = ProjectViewStatusNew
			ps.dialog.Show()
			ps.dialog.Reset()
			ps.pdb.Reset()
			ps.pdb.Focus()
		}
	}

	cmds = append(cmds, ps.dialog.Update(msg))
	if ps.dialog.IsDone() && ps.dialog.IsVisible() {
		if ps.dialog.GetConfirm() == "yes" {
			value := ps.dialog.Value
			if value != nil {
				if pjName, ok := value.(string); ok && pjName != "" {
					ps.AddProject(models.NewProject(uuid.NewString(), pjName))
				}
			}
		}
		ps.dialog.Hide()
		ps.viewstatus = ProjectViewStatusList
	}

	return tea.Batch(cmds...)
}

func (ps *ProjectsScreen) View() string {
	var content string

	if ps.viewstatus == ProjectViewStatusNew {
		content += ps.dialog.View()
		return content
	}

	if ps.viewstatus == ProjectViewStatusList {
		content = fmt.Sprintf("Projects (%d)\n\n", len(ps.collection))

		for _, project := range ps.collection {
			content += project.View() + "\n"
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, content)
	}

	return ""
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

func (s *ProjectsScreen) Display() string {
	count := len(s.collection)
	return fmt.Sprintf("Projects (%d)", count)
}

type ProjectDialogBody struct {
	name  string
	input textinput.Model
	done  bool
}

func (pdb *ProjectDialogBody) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	tiM, cmd := pdb.input.Update(msg)
	pdb.input = tiM
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			pdb.input.Blur()
			pdb.name = pdb.input.Value()
			pdb.done = true
		case "tab":
			pdb.input.Blur()
			pdb.name = pdb.input.Value()
		case "shift+tab":
			if !pdb.input.Focused() {
				pdb.input.Focus()
			}
		case "esc":
			pdb.input.Blur()
			pdb.name = ""
			pdb.done = false
		}
	}

	return tea.Batch(cmds...)
}

func (pdb *ProjectDialogBody) View() string {
	return pdb.input.View()
}

func (pdb *ProjectDialogBody) Value() any {
	return pdb.name
}

func (pdb *ProjectDialogBody) Reset() {
	pdb.name = ""
	pdb.done = false
	pdb.input.Reset()
	pdb.input.Blur()
}

func (pdb *ProjectDialogBody) Focus() {
	pdb.input.Focus()
}

func (pdb *ProjectDialogBody) Focused() bool {
	return pdb.input.Focused()
}

func NewProjectDialogBody() *ProjectDialogBody {
	pi := textinput.New()
	pi.Prompt = "▶ "
	pi.Placeholder = "Project Name"
	pi.Focus()
	return &ProjectDialogBody{
		name:  "",
		input: pi,
		done:  false,
	}
}

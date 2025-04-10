package screens

import (
	"p1/pkg/models"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectsScreen struct {
	collection    []*models.Project
	selected      int
	focused       bool
	search        string
	searchFocused bool
	viewport      viewport.Model
	width         int
	height        int
	ready         bool
}

func NewProjectsScreen(renderer *lipgloss.Renderer) *Screen {
	screen := &ProjectsScreen{
		collection:    []*models.Project{},
		selected:      0,
		focused:       false,
		search:        "",
		searchFocused: false,
		ready:         false,
		width:         30,
		height:        0,
	}
	return NewScreen(renderer, screen)
}

func (ps *ProjectsScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		ps.height = msg.Height
		if !ps.ready {
			ps.viewport = viewport.New(ps.width, msg.Height)
			ps.viewport.Style = ps.viewport.Style.
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				MaxHeight(ps.height - lipgloss.Height(ps.View()))
			ps.ready = true
		} else {
			ps.viewport.Width = ps.width
			ps.viewport.Height = ps.height
		}

	}
	for _, project := range ps.collection {
		cmds = append(cmds, project.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (ps *ProjectsScreen) View() string {
	content := ""
	for _, project := range ps.collection {
		content += project.View() + "\n"
	}
	contentHeight := lipgloss.Height(content)
	if contentHeight < ps.height {
		fillerHeight := ps.height - contentHeight
		filler := strings.Repeat(" \n", fillerHeight)
		content += filler
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

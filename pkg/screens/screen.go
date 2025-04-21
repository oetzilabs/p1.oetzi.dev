package screens

import (
	"fmt"
	"p1/pkg/dialog"
	"p1/pkg/interfaces"
	"p1/pkg/models"
	"p1/pkg/tui/theme"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Screen struct {
	ready     bool
	Content   interfaces.ScreenContent
	viewport  viewport.Model
	menuWidth int
	width     int
	height    int
	theme     theme.Theme
	focused   bool
	Commands  []*interfaces.FooterCommand
	footer    *models.Footer
}

func (s *Screen) FooterHeight() int {
	return lipgloss.Height(s.footer.View())
}

func (s *Screen) IsFocused() bool {
	return s.focused
}

func (s *Screen) SetFocused(focused bool) *Screen {
	s.focused = focused
	return s
}

func New(renderer *lipgloss.Renderer, content interfaces.ScreenContent, commands ...*interfaces.FooterCommand) *Screen {
	cmds := []*interfaces.FooterCommand{}
	cmds = append(cmds, commands...)
	theme := theme.BasicTheme(renderer, nil)
	footer := models.NewFooter(&theme, cmds)
	return &Screen{
		Content:  content,
		ready:    false,
		theme:    theme,
		focused:  false,
		Commands: cmds,
		footer:   footer,
	}
}

var modifiedKeyMap = viewport.KeyMap{
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "½ page down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	),
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}
	parentMsg := msg

	s.footer.SetCommands(s.Commands)

	cmds = append(cmds, s.footer.Update(parentMsg))

	footerHeight := lipgloss.Height(s.footer.View())

	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		s.height = msg.Height
		s.menuWidth = msg.MenuWidth
		s.width = msg.Width
		if !s.ready {
			s.viewport = viewport.New(msg.Width-msg.MenuWidth-2, msg.Height-footerHeight)
			s.viewport.HighPerformanceRendering = false
			s.viewport.KeyMap = modifiedKeyMap
			s.viewport.Style = s.viewport.Style.
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				PaddingBottom(1)
			s.ready = true
		} else {
			s.viewport.Width = max(0, msg.Width-msg.MenuWidth-2)
			s.viewport.Height = max(0, msg.Height-footerHeight)
			s.viewport.Style = s.viewport.Style.MaxHeight(msg.Height - footerHeight)
		}
	}

	s.viewport.KeyMap = modifiedKeyMap
	s.viewport.Width = max(0, s.width-s.menuWidth-2)
	s.viewport.Height = max(0, s.height-footerHeight)
	s.viewport.Style = s.viewport.Style.MaxHeight(max(0, s.height-footerHeight))
	cmds = append(cmds, s.Content.Update(parentMsg))
	s.viewport.SetContent(s.Content.View())

	if s.focused {
		vpm, cmd := s.viewport.Update(parentMsg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		s.viewport = vpm
	}

	return tea.Batch(cmds...)
}

func (s *Screen) View() string {
	viewportView := s.viewport.View()
	content := s.Content.View()
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			viewportView,
			s.getScrollbar(content),
		),
		s.footer.View())
}

func (s *Screen) getScrollbar(content string) string {
	y := s.viewport.YOffset
	vh := s.viewport.Height
	ch := lipgloss.Height(content)
	if vh >= ch {
		return ""
	}

	height := (vh * vh) / ch
	maxScroll := ch - vh
	nYP := 1.0 - (float64(y) / float64(maxScroll))
	if nYP <= 0 {
		nYP = 1
	} else if nYP >= 1 {
		nYP = 0
	}

	bar := s.theme.Base().
		Height(height).
		Width(1).
		Background(s.theme.Accent()).
		Render()

	style := s.theme.Base().Width(1).Height(vh)

	return style.Render(lipgloss.PlaceVertical(vh, lipgloss.Position(nYP), bar))
}

func (s *Screen) Display() string {
	if s.Content != nil {
		return s.Content.Display()
	}
	return ""
}

// Project Screen

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
	return New(renderer, screen, newProjectKey)
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

// Brokers Screen
type BrokersScreen struct {
	collection []*models.Broker
	selected   int
}

func NewBrokersScreen(renderer *lipgloss.Renderer) *Screen {
	screen := &BrokersScreen{
		collection: []*models.Broker{},
		selected:   0,
	}
	return New(renderer, screen, &interfaces.FooterCommand{Key: "n", Value: "New Broker"})
}

func (bs *BrokersScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for _, broker := range bs.collection {
		cmds = append(cmds, broker.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (bs *BrokersScreen) View() string {
	content := fmt.Sprintf("Brokers (%d)\n", len(bs.collection))
	for _, broker := range bs.collection {
		content += broker.View() + "\n"
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, content)
}

func (s *BrokersScreen) AddBroker(broker *models.Broker) *BrokersScreen {
	s.collection = append(s.collection, broker)
	return s
}

func (s *BrokersScreen) RemoveBroker(broker *models.Broker) *BrokersScreen {
	for i, p := range s.collection {
		if p.ID == broker.ID {
			s.collection = slices.Delete(s.collection, i, i+1)
			break
		}
	}
	return s
}

func (s *BrokersScreen) Display() string {
	count := len(s.collection)
	return fmt.Sprintf("Brokers (%d)", count)
}

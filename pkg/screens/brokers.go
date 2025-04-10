package screens

import (
	"p1/pkg/models"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BrokersScreen struct {
	collection    []*models.Broker
	selected      int
	focused       bool
	search        string
	searchFocused bool
	viewport      viewport.Model
	width         int
	height        int
	ready         bool
}

func NewBrokersScreen(renderer *lipgloss.Renderer) *Screen {
	screen := &BrokersScreen{
		collection:    []*models.Broker{},
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

func (bs *BrokersScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case models.InternalWindowSizeMsg:
		bs.height = msg.Height
		if !bs.ready {
			bs.viewport = viewport.New(bs.width, msg.Height)
			bs.viewport.Style = bs.viewport.Style.
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				MaxHeight(bs.height - lipgloss.Height(bs.View()))
			bs.ready = true
		} else {
			bs.viewport.Width = bs.width
			bs.viewport.Height = bs.height
		}

	}
	for _, broker := range bs.collection {
		cmds = append(cmds, broker.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (bs *BrokersScreen) View() string {
	content := ""
	for _, broker := range bs.collection {
		content += broker.View() + "\n"
	}
	contentHeight := lipgloss.Height(content)
	if contentHeight < bs.height {
		fillerHeight := bs.height - contentHeight
		filler := strings.Repeat(" \n", fillerHeight)
		content += filler
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

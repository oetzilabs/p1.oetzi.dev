package tui

import (
	"fmt"
	models "p1/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServerCollection struct {
	servers  []*models.Server
	selected int
}

func NewServerCollection() *ServerCollection {
	return &ServerCollection{
		servers: []*models.Server{},
	}
}

func (sc *ServerCollection) AddServer(server *models.Server) {
	sc.servers = append(sc.servers, server)
}

func (sc *ServerCollection) SelectServer(id string) {
	for i, server := range sc.servers {
		if server.ID == id {
			sc.selected = i
			return
		}
	}
}

func (sc *ServerCollection) RemoveServer(id string) {
	for i, server := range sc.servers {
		if server.ID == id {
			sc.servers = append(sc.servers[:i], sc.servers[i+1:]...)
			return
		}
	}
}

func (sc *ServerCollection) Update(msg tea.Msg) tea.Cmd {
	parentMsg := msg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		}
	}
	var cmd tea.Cmd

	if len(sc.servers) > 0 {
		cmd = tea.Batch(sc.servers[sc.selected].Update(parentMsg), cmd)
	}

	return cmd
}

func (sc *ServerCollection) View() string {
	if len(sc.servers) == 0 {
		return "No servers available. Press 'n' to add a new server."
	}
	content := sc.servers[sc.selected].View()

	return content
}

func (sc *ServerCollection) Display() string {
	style := lipgloss.NewStyle()
	return lipgloss.JoinHorizontal(lipgloss.Left, style.Render("Servers", "__FILLER__", fmt.Sprintf("(%d)", len(sc.servers))))
}

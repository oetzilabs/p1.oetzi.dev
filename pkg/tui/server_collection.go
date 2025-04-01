package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerCollection struct {
	servers  []*ServerView
	selected int
}

func NewServerCollection() *ServerCollection {
	return &ServerCollection{
		servers: []*ServerView{},
	}
}

func (sc *ServerCollection) AddServer(server *ServerView) {
	sc.servers = append(sc.servers, server)
}

func (sc *ServerCollection) SelectServer(id string) {
	for i, server := range sc.servers {
		if server.Id == id {
			sc.selected = i
			return
		}
	}
}

func (sc *ServerCollection) RemoveServer(id string) {
	for i, server := range sc.servers {
		if server.Id == id {
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
	if len(sc.servers) == 0 {
		return "Servers"
	}
	return fmt.Sprintf("Servers (%d)", len(sc.servers))
}

package tui

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerCollection struct {
	servers        []*ServerView
	selected       int
	to_remove      string
	confirm_delete *Confirmation
	dialog         *Dialog
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

func (sc *ServerCollection) ConfirmRemoveProject(id string) {
	sc.to_remove = id
	name := sc.servers[sc.selected].Name
	confirm := NewConfirmation("Do you really wish to delete"+name+"?", func() {
		sc.RemoveServer(sc.to_remove)
	}, func() {
		sc.to_remove = ""
	})
	sc.confirm_delete = confirm
}

func (sc *ServerCollection) AddServerDialog() {
	inputs := []Input{
		*NewInput("ID", "123", nil),
		*NewInput("Name", "My Project", nil),
		*NewInput("URL", "https://example.com", nil),
	}
	dialog := NewDialog(
		"Enter the name of the project",
		inputs,
		func(values interface{}) {
			v := reflect.ValueOf(values)
			id := v.FieldByName("ID").String()
			name := v.FieldByName("Name").String()
			url := v.FieldByName("URL").String()
			sc.AddServer(NewServer(id, name, url))
			sc.dialog = nil
		},
		func() {
			sc.dialog = nil
		},
	)
	sc.dialog = dialog
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
	if sc.confirm_delete != nil {
		cmd = sc.confirm_delete.Update(msg)
	}

	if len(sc.servers) > 0 {
		cmd = tea.Batch(sc.servers[sc.selected].Update(parentMsg), cmd)
	}

	// Send display information up to parent
	updateCmd := tea.Cmd(func() tea.Msg {
		return UpdateTabDisplay{
			DisplayLeft:  "Servers",
			DisplayRight: []string{fmt.Sprintf("(%d)", len(sc.servers))},
		}
	})
	cmd = tea.Batch(cmd, updateCmd)
	return cmd
}

func (sc *ServerCollection) View() string {
	content := sc.servers[sc.selected].View()
	if sc.to_remove != "" {
		content += sc.confirm_delete.View()
	}
	return content
}

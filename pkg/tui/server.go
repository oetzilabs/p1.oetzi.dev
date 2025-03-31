package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServerView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func NewServer(id string, name string, url string) *ServerView {
	return &ServerView{
		Id:   id,
		Name: name,
		Url:  url,
	}
}

func (s *ServerView) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (s *ServerView) View() string {
	mainStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#212121")).
		Padding(0, 1)

	return mainStyle.Render(s.Name)
}

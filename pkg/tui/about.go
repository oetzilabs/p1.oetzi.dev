package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type About struct {
}

func NewAbout() *About {
	return &About{}
}

func (a *About) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (a *About) View() string {
	return "About Content"
}

func (a *About) Display() string {
	return "About"
}

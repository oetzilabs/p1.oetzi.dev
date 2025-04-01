package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type EmptyContent struct {
}

func NewEmptyContent() *EmptyContent {
	return &EmptyContent{}
}

func (e *EmptyContent) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to EmptyContent
	return nil
}

func (e *EmptyContent) View() string {
	var content string = ""
	return content
}

func (e *EmptyContent) Display() string {
	return ""
}

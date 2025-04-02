package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ExitTab struct{}

func NewExitTab() Tab {
	return Tab{
		ID:      "exit",
		Hidden:  false,
		Group:   AlignBottom,
		Content: &ExitTab{},
	}
}

func (et *ExitTab) Update(msg tea.Msg) tea.Cmd {
	// if clicked enter, quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return tea.Quit
		}
	}
	return nil
}

func (et *ExitTab) View() string {
	return "This is the Exit screen. Press 'ctrl+c' to quit."
}

func (et *ExitTab) Display() string {
	return "Exit"
}

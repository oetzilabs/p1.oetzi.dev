package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ExitTab struct{}

func NewExitTab() Tab {
	return Tab{
		ID:      "exit",
		Hidden:  false,
		Group:   AlignBottom,
		Content: &ExitTab{},
		Helper:  "Press 'q' to quit.",
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
	return ""
}

func (et *ExitTab) Display() string {
	style := lipgloss.NewStyle()
	return lipgloss.JoinVertical(lipgloss.Left, style.Render("Exit", "__FILLER__", "(q)"))
}

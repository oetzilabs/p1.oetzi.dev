package tabs

import tea "github.com/charmbracelet/bubbletea"

type AboutTab struct{}

func NewAboutTab() Tab {
	return Tab{
		ID:      "about",
		Group:   AlignBottom,
		Hidden:  false,
		Content: &AboutTab{},
		Helper:  "",
	}
}

func (at *AboutTab) Update(msg tea.Msg) tea.Cmd {
	// No updates needed for the About tab
	return nil
}

func (at *AboutTab) View() string {
	return "This is the About screen. Press 'ctrl+c' to quit."
}

func (at *AboutTab) Display() string {
	return "About"
}

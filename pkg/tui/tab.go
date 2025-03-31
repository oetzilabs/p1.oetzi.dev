package tui

type Tab struct {
	ID      screen
	Name    string
	Hidden  bool
	Content Content
}

func NewTab(id screen, name string, hidden bool, content Content) Tab {
	return Tab{
		ID:      id,
		Name:    name,
		Hidden:  hidden,
		Content: content,
	}
}

package tui

func (m model) LogoView() string {
	return m.theme.TextAccent().Bold(true).Render("p1.oetzi.dev ") + m.CursorView()
}

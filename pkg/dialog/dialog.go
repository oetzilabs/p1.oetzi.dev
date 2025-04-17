package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DialogBody interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
	Value() any
	Focused() bool
}

type Dialog struct {
	title         string
	body          DialogBody
	Value         any
	done          bool
	width         int
	height        int
	confirmStatus string
	visible       bool
}

var (
	dialogBoxStyle = lipgloss.NewStyle().MarginRight(1)
	buttonStyle    = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3)
	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Underline(true)
	subtle                = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	dialogBackgroundColor = lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#222222"}
)

func NewDialog(title string, body DialogBody) *Dialog {
	return &Dialog{
		width:         0,
		height:        0,
		title:         title,
		body:          body,
		done:          false,
		confirmStatus: "yes",
		visible:       false,
	}
}

func (d *Dialog) UpdateSize(w, h int) {
	d.width = w
	d.height = h
}

func (d *Dialog) View() string {
	if !d.visible {
		return ""
	}
	if d.done {
		return ""
	}

	question := lipgloss.NewStyle().Align(lipgloss.Left).Render(d.title + "\n")
	body := d.body.View() + "\n"

	var okButton string
	cancelButton := buttonStyle.Render("Cancel")
	if d.confirmStatus == "yes" {
		okButton = activeButtonStyle.Render("Yes")
	} else {
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("Cancel")
	}

	// Create a filler to push buttons to the right
	halfWidth := d.width / 2
	fillerWidth := max(0, halfWidth-lipgloss.Width(cancelButton)-lipgloss.Width(okButton)-4)

	filler := lipgloss.NewStyle().Width(fillerWidth).Render("")

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, filler, cancelButton,
		lipgloss.NewStyle().Background(dialogBackgroundColor).Render(" "),
		okButton)

	content := lipgloss.JoinVertical(lipgloss.Top, question, body, buttons)

	box := lipgloss.NewStyle().
		Background(dialogBackgroundColor).
		Padding(1).
		Render(content)

	dialog := lipgloss.Place(d.width, d.height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(box),
		lipgloss.WithWhitespaceChars("猫咪"),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	return dialog
}

func (d *Dialog) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}
	if d.done {
		return tea.Batch(cmds...)
	}
	if !d.visible {
		return tea.Batch(cmds...)
	}

	cmds = append(cmds, d.body.Update(msg))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !d.done {
				d.Value = d.body.Value()
				d.done = true
			}
		case "esc":
			d.Value = nil
			d.done = true
		case "left":
			if !d.body.Focused() {
				d.confirmStatus = "cancel"
			}
		case "right":
			if !d.body.Focused() {
				d.confirmStatus = "yes"
			}
		}
	}

	return tea.Batch(cmds...)
}

func (d *Dialog) GetConfirm() string {
	return d.confirmStatus
}

func (d *Dialog) IsDone() bool {
	return d.done
}

func (d *Dialog) IsVisible() bool {
	return d.visible
}

func (d *Dialog) Hide() {
	d.visible = false
}

func (d *Dialog) Show() {
	d.visible = true
}

func (d *Dialog) Reset() {
	d.done = false
	d.confirmStatus = "yes"
}

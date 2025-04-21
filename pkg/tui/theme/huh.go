package theme

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
)

// copy returns a copy of a TextInputStyles with all children styles copied.
func copyTextStyles(t huh.TextInputStyles) huh.TextInputStyles {
	return huh.TextInputStyles{
		Cursor:      t.Cursor,
		Placeholder: t.Placeholder,
		Prompt:      t.Prompt,
		Text:        t.Text,
	}
}

// copy returns a copy of a FieldStyles with all children styles copied.
func copyFieldStyles(f huh.FieldStyles) huh.FieldStyles {
	return huh.FieldStyles{
		Base:           f.Base,
		Title:          f.Title,
		Description:    f.Description,
		ErrorIndicator: f.ErrorIndicator,
		ErrorMessage:   f.ErrorMessage,
		SelectSelector: f.SelectSelector,
		// NextIndicator:       f.NextIndicator,
		// PrevIndicator:       f.PrevIndicator,
		Option: f.Option,
		// Directory:           f.Directory,
		// File:                f.File,
		MultiSelectSelector: f.MultiSelectSelector,
		SelectedOption:      f.SelectedOption,
		SelectedPrefix:      f.SelectedPrefix,
		UnselectedOption:    f.UnselectedOption,
		UnselectedPrefix:    f.UnselectedPrefix,
		FocusedButton:       f.FocusedButton,
		BlurredButton:       f.BlurredButton,
		TextInput:           copyTextStyles(f.TextInput),
		Card:                f.Card,
		NoteTitle:           f.NoteTitle,
		Next:                f.Next,
	}
}

func copy(t huh.Theme) huh.Theme {
	return huh.Theme{
		Form:           t.Form,
		Group:          t.Group,
		FieldSeparator: t.FieldSeparator,
		Blurred:        copyFieldStyles(t.Blurred),
		Focused:        copyFieldStyles(t.Focused),
		Help: help.Styles{
			Ellipsis:       t.Help.Ellipsis,
			ShortKey:       t.Help.ShortKey,
			ShortDesc:      t.Help.ShortDesc,
			ShortSeparator: t.Help.ShortSeparator,
			FullKey:        t.Help.FullKey,
			FullDesc:       t.Help.FullDesc,
			FullSeparator:  t.Help.FullSeparator,
		},
	}
}

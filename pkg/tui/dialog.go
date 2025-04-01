package tui

import (
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dialog struct {
	message   string
	inputs    []Input
	onConfirm func(interface{})
	onCancel  func()
}

// NewDialog creates a new dialog with named inputs
// The onConfirm callback will receive a struct with fields matching the input names
// Example:
//
//	inputs := []Input{
//	  *NewInput("Name", "", InputOptions{}),
//	  *NewInput("URL", "", InputOptions{}),
//	}
//	dialog := NewDialog("Enter details", inputs, func(values interface{}) {
//	  v := reflect.ValueOf(values)
//	  name := v.FieldByName("Name").String()
//	  url := v.FieldByName("URL").String()
//	}, func() {})
func NewDialog(message string, inputs []Input, onConfirm func(interface{}), onCancel func()) *Dialog {
	return &Dialog{
		message:   message,
		inputs:    inputs,
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
}

// IsValid checks if all inputs have non-empty values
func (d *Dialog) IsValid() bool {
	if len(d.inputs) == 0 {
		return false
	}
	for _, input := range d.inputs {
		if input.Value == "" {
			return false
		}
	}
	return true
}

func (d *Dialog) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			if d.IsValid() {
				// Create a struct with fields matching input names
				fields := make([]reflect.StructField, len(d.inputs))
				values := make([]reflect.Value, len(d.inputs))

				for i, input := range d.inputs {
					// Convert input name to PascalCase for struct field
					fieldName := strings.Title(strings.ToLower(input.Label))
					fields[i] = reflect.StructField{
						Name: fieldName,
						Type: reflect.TypeOf(""),
					}
					values[i] = reflect.ValueOf(input.Value)
				}

				// Create the struct type and value
				structType := reflect.StructOf(fields)
				structValue := reflect.New(structType).Elem()

				// Set the values
				for i, value := range values {
					structValue.Field(i).Set(value)
				}

				d.onConfirm(structValue.Interface())
			}
		case "n":
			d.onCancel()
		}
	}
	return nil
}

func (d *Dialog) View() string {
	inputs := make([]string, len(d.inputs))
	for i, input := range d.inputs {
		inputs[i] = input.View()
	}

	validStatus := ""
	if !d.IsValid() {
		validStatus = " (Please fill in all fields)"
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		d.message,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			inputs...,
		),
		"(y) Yes  |  (n) No"+validStatus,
	)
}

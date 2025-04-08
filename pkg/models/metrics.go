package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Metrics struct {
	CPU     float64 `json:"cpu"`
	RAM     float64 `json:"ram"`
	Disk    float64 `json:"disk"`
	Network float64 `json:"network"`
}

func (m *Metrics) Update(msg tea.Msg) tea.Cmd {
	// Handle updates specific to Metrics
	return nil
}

func (m *Metrics) View() string {
	mainStyle := lipgloss.NewStyle().Padding(2)
	content := fmt.Sprintf("CPU: %.2f%%\nRAM: %.2f%%\nDisk: %.2f%%\nNetwork: %.2f%%", m.CPU, m.RAM, m.Disk, m.Network)
	return mainStyle.Render(content)
}

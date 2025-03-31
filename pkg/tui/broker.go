package tui

import (
	"p1/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
)

type Broker struct {
	Servers []models.Server `json:"servers"`
}

type BrokerData struct {
	Data Broker `json:"data"`
}

type BrokerView struct {
	Id   string
	Name string
	URL  string
}

func NewBroker(id string, name string, url string) *BrokerView {
	return &BrokerView{
		Id:   id,
		Name: name,
		URL:  url,
	}
}

func (b *BrokerView) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (b *BrokerView) View() string {
	return "Broker: " + b.Name + " (" + b.URL + ")\n"
}

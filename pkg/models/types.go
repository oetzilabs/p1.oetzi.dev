package models

// State represents the current state of the system
type State struct {
	Servers  []Server  `json:"servers"`
	Brokers  []Broker  `json:"brokers"`
	Projects []Project `json:"projects"`
	Actor    Actor     `json:"actor"`
}

// Update represents a message sent over the websocket
type Update struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type VisibleError struct {
	Message string
}

func NewVisibleError(message string) VisibleError {
	return VisibleError{
		Message: message,
	}
}

type InternalWindowSizeMsg struct {
	Width        int
	Height       int
	MenuWidth    int
	FooterHeight int
}

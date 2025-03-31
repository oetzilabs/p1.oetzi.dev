package models

import "github.com/google/uuid"

// State represents the current state of the system
type State struct {
	Servers  []Server  `json:"servers"`
	Brokers  []Broker  `json:"brokers"`
	Projects []Project `json:"projects"`
	Actor    Actor     `json:"actor"`
}

type Server struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Broker struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Actor struct {
	ID string `json:"id"`
}

// Update represents a message sent over the websocket
type Update struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewActor() Actor {
	id := uuid.New().String()
	return Actor{
		ID: id,
	}
}

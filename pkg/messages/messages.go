package messages

import (
	"p1/pkg/states"
)

type RerenderMessage struct {
	Key      string
	Value    any
	OldValue any
}

type SyncMsg *states.ClientState

type MessageType string

const (
	TypeListServices    MessageType = "LIST_SERVICES"
	TypeRegisterService MessageType = "REGISTER_SERVICE"
	TypeRemoveService   MessageType = "REMOVE_SERVICE"

	TypeListBrokers    MessageType = "LIST_BROKERS"
	TypeRegisterBroker MessageType = "REGISTER_BROKER"
	TypeRemoveBroker   MessageType = "REMOVE_BROKER"

	TypeListProjects     MessageType = "LIST_PROJECTS"
	TypeRegisterProjects MessageType = "REGISTER_PROJECTS"
	TypeRemoveProjects   MessageType = "REMOVE_PROJECTS"

	TypeMetrics   MessageType = "METRICS"
	TypeBroadcast MessageType = "BROADCAST"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
	Sender  string      `json:"sender"`
}

package client

import (
	"fmt"
	"log/slog"
	"p1/pkg/server"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	link     string
	status   WebsocketServerStatus
	services []*server.Service
	brokers  []string
	Conn     *websocket.Conn
}

type WebsocketServerStatus string

const (
	WebsocketStatusOK                 WebsocketServerStatus = "OK"
	WebsocketStatusWarning            WebsocketServerStatus = "WARNING"
	WebsocketStatusCritical           WebsocketServerStatus = "CRITICAL"
	WebsocketStatusNotEnoughResources WebsocketServerStatus = "NOT_ENOUGH_RESOURCES"
	WebsocketStatusMissingLink        WebsocketServerStatus = "MISSING_LINK"
	WebsocketStatusUnknown            WebsocketServerStatus = "UNKNOWN"
)

func (ws *WebsocketClient) Update() error {
	if ws.Conn == nil {
		return fmt.Errorf("no websocket available")
	}
	// if ws.Conn == nil {
	// 	// first create the connection
	// 	dialer := websocket.Dialer{
	// 		HandshakeTimeout: DEFAULT_TIMEOUT,
	// 	}
	// 	connection, _, err := dialer.Dial(ws.link, nil)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	ws.Conn = connection
	// 	// defer connection.Close()
	// }
	// Send request for Services
	msg := server.Message{
		Type:    server.TypeListServices,
		Payload: nil,
	}
	if err := ws.Conn.WriteJSON(msg); err != nil {
		return err
	}

	if err := ws.Conn.ReadJSON(&msg); err != nil {
		return err
	}

	ws.services = msg.Payload.([]*server.Service)

	// Send request for Brokers
	msg = server.Message{
		Type:    server.TypeListBrokers,
		Payload: nil,
	}
	if err := ws.Conn.WriteJSON(msg); err != nil {
		return err
	}

	if err := ws.Conn.ReadJSON(&msg); err != nil {
		return err
	}

	ws.brokers = msg.Payload.([]string)

	return nil
}

func NewWebsocketClient(link string) *WebsocketClient {
	return &WebsocketClient{
		link:     link,
		services: []*server.Service{},
		brokers:  []string{},
		status:   WebsocketStatusUnknown,
	}
}

func (ws *WebsocketClient) CreateConnection() {
	dialer := websocket.Dialer{
		HandshakeTimeout: DEFAULT_TIMEOUT,
	}
	connection, _, err := dialer.Dial(ws.link, nil)
	if err != nil {
		slog.Error("Error creating websocket connection", "error", err.Error())
		return
	}
	ws.Conn = connection
}

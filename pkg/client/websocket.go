package client

import (
	"fmt"
	"net/http"
	"net/url"
	"p1/pkg/broker"
	"p1/pkg/server"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	ID       string
	link     string
	status   ServerStatus
	services []*server.Service
	brokers  []*broker.Broker
}

type ServerStatus string

const (
	StatusOK                 ServerStatus = "OK"
	StatusWarning            ServerStatus = "WARNING"
	StatusCritical           ServerStatus = "CRITICAL"
	StatusNotEnoughResources ServerStatus = "NOT_ENOUGH_RESOURCES"
	StatusMissingLink        ServerStatus = "MISSING_LINK"
	StatusUnknown            ServerStatus = "UNKNOWN"
)

func isAcceptableStatus(status ServerStatus) bool {
	switch status {
	case StatusOK:
		return true
	case StatusWarning:
		return true
	case StatusCritical:
		return false
	case StatusNotEnoughResources:
		return false
	case StatusMissingLink:
		return false
	default:
		return false
	}
}

func (ws *WebsocketServer) StatusCheck() ServerStatus {
	if ws.link == "" {
		return StatusMissingLink
	}

	// Parse the URL
	u, err := url.Parse(ws.link)
	if err != nil {
		return StatusCritical
	}

	// Try HTTP connection first
	client := http.Client{
		Timeout: DEFAULT_TIMEOUT,
	}
	_, err = client.Get(fmt.Sprintf("http://%s", u.Host))
	if err != nil {
		return StatusCritical
	}

	// Try WebSocket connection
	dialer := websocket.Dialer{
		HandshakeTimeout: DEFAULT_TIMEOUT,
	}
	c, _, err := dialer.Dial(ws.link, nil)
	if err != nil {
		return StatusWarning // Server is up but WebSocket failed
	}
	defer c.Close()

	return StatusOK
}

func (ws *WebsocketServer) Update() error {
	if ws.link == "" {
		return fmt.Errorf("link is empty")
	}

	// Parse the URL
	_, err := url.Parse(ws.link)
	if err != nil {
		return fmt.Errorf("invalid link: %v", err)
	}

	// Connect via temporary WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: DEFAULT_TIMEOUT,
	}
	c, _, err := dialer.Dial(ws.link, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	ws.status = ws.StatusCheck()

	// Send request for Services
	msg := server.Message{
		Type:    server.TypeListServices,
		Payload: nil,
	}
	if err := c.WriteJSON(msg); err != nil {
		return err
	}

	if err := c.ReadJSON(&msg); err != nil {
		return err
	}

	services := []*server.Service{}
	tempServices := msg.Payload.([]*server.Service)
	for _, svc := range tempServices {
		if !ws.hasService(svc.ID) {
			if svc.Endpoint != "" {
				services = append(services, svc)
			}
		}
	}
	ws.services = services

	// Send request for Services
	msg = server.Message{
		Type:    server.TypeListBrokers,
		Payload: nil,
	}
	if err := c.WriteJSON(msg); err != nil {
		return err
	}

	if err := c.ReadJSON(&msg); err != nil {
		return err
	}

	brokers := []*broker.Broker{}
	tempBrokers := msg.Payload.([]*broker.Broker)
	for _, brk := range tempBrokers {
		if !ws.hasBroker(brk.ID) {
			if brk.URL != "" {
				brokers = append(brokers, brk)
			}
		}
	}
	ws.brokers = brokers

	return nil
}

func (ws *WebsocketServer) hasService(id string) bool {
	for _, svc := range ws.services {
		if svc.ID == id {
			return true
		}
	}
	return false
}

func (ws *WebsocketServer) hasBroker(id string) bool {
	for _, brk := range ws.brokers {
		if brk.ID == id {
			return true
		}
	}
	return false
}

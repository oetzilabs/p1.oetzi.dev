package client

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"p1/pkg/server"
	"slices"
	"time"

	"github.com/gorilla/websocket"
)

const DEFAULT_TIMEOUT = 10 * time.Second
const DEFAULT_TICKER_INTERVAL = 5 * time.Second

type WebsocketServer struct {
	ID       string
	link     string
	status   ServerStatus
	services []*server.Service
}

type ServerStatus string

const (
	StatusOK                 ServerStatus = "OK"
	StatusWarning            ServerStatus = "WARNING"
	StatusCritical           ServerStatus = "CRITICAL"
	StatusNotEnoughResources ServerStatus = "NOT_ENOUGH_RESOURCES"
	StatusMissingLink        ServerStatus = "MISSING_LINK"
)

func isAcceptable(status ServerStatus) bool {
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

func (ws *WebsocketServer) UpdateServices() []*server.Service {
	if ws.link == "" {
		return []*server.Service{}
	}

	// Parse the URL
	_, err := url.Parse(ws.link)
	if err != nil {
		return []*server.Service{}
	}

	// Connect via temporary WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: DEFAULT_TIMEOUT,
	}
	c, _, err := dialer.Dial(ws.link, nil)
	if err != nil {
		return []*server.Service{}
	}
	defer c.Close()

	// Send request
	msg := server.Message{
		Type:    server.TypeListServices,
		Payload: nil,
	}
	if err := c.WriteJSON(msg); err != nil {
		return []*server.Service{}
	}

	if err := c.ReadJSON(&msg); err != nil {
		return []*server.Service{}
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
	return services
}

func (ws *WebsocketServer) hasService(id string) bool {
	for _, svc := range ws.services {
		if svc.ID == id {
			return true
		}
	}
	return false
}

type Client struct {
	wsServers []*WebsocketServer
	ticker    *time.Ticker
}

func NewClient() *Client {
	return &Client{
		wsServers: []*WebsocketServer{},
		ticker:    time.NewTicker(DEFAULT_TICKER_INTERVAL),
	}
}

func (c *Client) AddServer(id string, link string) error {
	server := WebsocketServer{
		ID:   id,
		link: link,
	}

	server.status = server.StatusCheck()

	if !isAcceptable(server.status) {
		return fmt.Errorf("server %s can not be added, status: %v", id, server.status)
	} else {
		c.wsServers = append(c.wsServers, &server)
	}
	return nil
}

func (c *Client) RemoveServer(id string) error {
	for i, server := range c.wsServers {
		if server.ID == id {
			c.wsServers = slices.Delete(c.wsServers, i, i+1)
			return nil
		}
	}
	// server not found
	return fmt.Errorf("server %s not found", id)
}

func (c *Client) GetServers() []*WebsocketServer {
	return c.wsServers
}

func (c *Client) GetServer(id string) *WebsocketServer {
	for _, server := range c.wsServers {
		if server.ID == id {
			return server
		}
	}
	// server not found
	return nil
}

func (c *Client) UpdateServer(id string) error {
	slog.Info("Updating server", "id", id)
	server := c.GetServer(id)
	if server == nil {
		return fmt.Errorf("server %s not found", id)
	}
	server.status = server.StatusCheck()
	server.services = server.UpdateServices()
	return nil
}

func (c *Client) Sync() error {
	slog.Info("Syncing servers")
	if len(c.wsServers) == 0 {
		return fmt.Errorf("no servers available")
	}
	for _, server := range c.wsServers {
		err := c.UpdateServer(server.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Start() {
	// make a ticker that check syncs every 5 seconds.
	go func() {
		for range c.ticker.C {
			c.Sync()
		}
	}()
}

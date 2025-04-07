package client

import (
	"fmt"
	"log/slog"
	"slices"
	"time"
)

const DEFAULT_TIMEOUT = 10 * time.Second
const DEFAULT_TICKER_INTERVAL = 5 * time.Second

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
		ID:     id,
		link:   link,
		status: StatusUnknown,
	}

	server.status = server.StatusCheck()

	if !isAcceptableStatus(server.status) {
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
	server.Update()
	return nil
}

func (c *Client) Sync() error {
	slog.Info("Syncing servers")
	if len(c.wsServers) == 0 {
		return nil
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

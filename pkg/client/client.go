package client

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_TIMEOUT = 10 * time.Second

type Client struct {
	WsClient *WebsocketClient
	cid      string
}

func NewClient() *Client {
	cid := uuid.New().String()
	client := NewWebsocketClient(fmt.Sprintf("ws://localhost:28080/ws?CID=%s", cid))

	return &Client{
		cid:      cid,
		WsClient: client,
	}
}

func (c *Client) Sync() error {
	slog.Info("Syncing servers")
	err := c.WsClient.Update()
	if err != nil {
		return err
	}
	// TODO: update servers + brokers

	return nil
}

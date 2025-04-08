package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"p1/pkg/messages"
	"p1/pkg/models"
	"p1/pkg/states"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	DEFAULT_TIMEOUT         = 10 * time.Second
	MAX_RECONNECT_ATTEMPTS  = 5
	INITIAL_RECONNECT_DELAY = 1 * time.Second
)

type Client struct {
	cid   string
	link  string
	state *states.ClientState
	conn  *websocket.Conn
}

func NewClient(mainServerLink string) *Client {
	cid := uuid.New().String()
	c := &Client{
		cid:  cid,
		link: mainServerLink,
		state: &states.ClientState{
			Projects: []*models.Project{},
			Servers:  []*models.Server{},
			Brokers:  []*models.Broker{},
		},
	}
	return c
}

func (c *Client) Init() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: DEFAULT_TIMEOUT,
	}

	headers := http.Header{}
	headers.Add("X-Client-Id", c.cid)

	conn, resp, err := dialer.Dial(c.link, headers)
	if err != nil {
		return fmt.Errorf("failed to Dial %v", err.Error())
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	c.conn = conn
	return nil
}

func (c *Client) Stop() error {
	if c.conn == nil {
		return fmt.Errorf("there is no active connection")
	}
	return c.conn.Close()
}

func (c *Client) reconnect() error {
	delay := INITIAL_RECONNECT_DELAY
	attempts := 0

	for attempts < MAX_RECONNECT_ATTEMPTS {
		slog.Info("attempting to reconnect", "attempt", attempts+1)
		err := c.Init()
		if err == nil {
			slog.Info("reconnected successfully")
			return nil
		}

		attempts++
		slog.Error("reconnection failed", "error", err, "attempt", attempts)
		time.Sleep(delay)
		delay *= 2 // exponential backoff
	}

	return fmt.Errorf("failed to reconnect after %d attempts", MAX_RECONNECT_ATTEMPTS)
}

func (c *Client) processMessage(msg *messages.Message) error {
	switch msg.Type {
	case messages.TypeMetrics:
		if metrics, ok := msg.Payload.(*models.Metrics); ok {
			c.state.Metrics = metrics
		}
	case messages.TypeListServices:
		if servers, ok := msg.Payload.([]*models.Server); ok {
			c.state.Servers = servers
		}
	case messages.TypeListBrokers:
		if brokers, ok := msg.Payload.([]*models.Broker); ok {
			c.state.Brokers = brokers
		}
	}
	return nil
}

func (c *Client) Start() error {
	if c.conn == nil {
		return fmt.Errorf("no active connection")
	}

	go func() {
		for {
			messageType, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("websocket error", "error", err)
					if reconnectErr := c.reconnect(); reconnectErr != nil {
						slog.Error("failed to reconnect", "error", reconnectErr)
						return
					}
					continue
				}
				return
			}

			switch messageType {
			case websocket.TextMessage:
				slog.Info("received text message", "message", string(message))
				var msg messages.Message
				if err := json.Unmarshal(message, &msg); err != nil {
					slog.Error("failed to unmarshal message", "error", err)
					continue
				}

				if err := c.processMessage(&msg); err != nil {
					slog.Error("failed to process message", "error", err)
					continue
				}
			case websocket.BinaryMessage:
				slog.Info("received binary message", "size", len(message))

			}
		}
	}()

	return nil
}

func (c *Client) Pull() *states.ClientState {
	return c.state
}

func (c *Client) SendMessage(msg *messages.Message) error {
	return c.conn.WriteJSON(msg)
}

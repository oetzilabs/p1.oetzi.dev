package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

// import "time"

type BrokerState struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type WebSocketMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WebSocketClient struct {
	conn     *websocket.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	url      string
	handlers map[string]func([]byte)
	msgChan  chan WebSocketMessage
}

func NewWebSocketClient(url string) *WebSocketClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketClient{
		url:      url,
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]func([]byte)),
		msgChan:  make(chan WebSocketMessage, 100),
	}
}

func (c *WebSocketClient) Connect() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *WebSocketClient) Disconnect() {
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.msgChan)
}

func (c *WebSocketClient) Send(data interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("websocket not connected")
	}
	return c.conn.WriteJSON(data)
}

func (c *WebSocketClient) Subscribe() error {
	if c.conn == nil {
		return fmt.Errorf("websocket not connected")
	}

	go func() {
		defer c.Disconnect()

		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					slog.Error("Error reading websocket message", "error", err)
					return
				}

				var msg WebSocketMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					slog.Error("Error unmarshaling message", "error", err)
					continue
				}

				// Send message to channel
				select {
				case c.msgChan <- msg:
				default:
					slog.Warn("Message channel full, dropping message")
				}

				// Call handler if registered
				if handler, ok := c.handlers[msg.Type]; ok {
					handler(msg.Data)
				}
			}
		}
	}()

	return nil
}

func (c *WebSocketClient) OnMessage(msgType string, handler func([]byte)) {
	c.handlers[msgType] = handler
}

func (c *WebSocketClient) MessageChannel() <-chan WebSocketMessage {
	return c.msgChan
}

func FetchBrokerState() (*BrokerState, error) {
	// TODO: Implement actual API call to fetch broker state
	// For now, return mock data
	time.Sleep(time.Second) // Simulate network delay
	return &BrokerState{
		Servers: []Server{
			{
				ID:   "1",
				Name: "Test Server",
				URL:  "http://localhost:8080",
			},
		},
	}, nil
}

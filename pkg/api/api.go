package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

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
type WebSocketStatus = string

const (
	WebSocketConnected    WebSocketStatus = "websocket:connected"
	WebSocketDisconnected WebSocketStatus = "websocket:disconnected"
	WebSocketUnknown      WebSocketStatus = "websocket:unknown"
)

type WebSocketApiClient struct {
	conn     *websocket.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	url      string
	handlers map[string]func([]byte)
	msgChan  chan WebSocketMessage
	status   WebSocketStatus
	mu       sync.Mutex
	once     sync.Once
	done     chan struct{} // new done channel
}

func NewWebSocketApiClient(url string) *WebSocketApiClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketApiClient{
		url:      url,
		ctx:      ctx,
		cancel:   cancel,
		handlers: make(map[string]func([]byte)),
		msgChan:  make(chan WebSocketMessage, 100),
		status:   WebSocketUnknown,
		mu:       sync.Mutex{},
		once:     sync.Once{},
		done:     make(chan struct{}), // initialize done channel
	}
}

func (c *WebSocketApiClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	slog.Info("Connecting to WebSocket", "url", c.url)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		slog.Error("Failed to connect to websocket", "error", err, "url", c.url)
		c.status = WebSocketUnknown
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	c.conn = conn
	c.status = WebSocketConnected
	slog.Info("Connected to WebSocket", "url", c.url)
	return nil
}

func (c *WebSocketApiClient) Disconnect() {
	c.once.Do(func() {
		slog.Info("Disconnecting WebSocket", "url", c.url)
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.cancel != nil {
			slog.Info("Cancelling context", "url", c.url)
			c.cancel()
		}

		// Signal Subscribe goroutine to exit
		slog.Info("Closing done channel", "url", c.url)
		close(c.done)
		slog.Info("Done channel closed", "url", c.url)

		if c.conn != nil {
			slog.Info("Closing WebSocket connection", "url", c.url)
			if err := c.conn.Close(); err != nil {
				slog.Error("Error closing WebSocket connection", "error", err, "url", c.url)
			}
			c.conn = nil // prevent double close
			slog.Info("WebSocket connection closed", "url", c.url)
		} else {
			slog.Warn("WebSocket connection already nil, skipping close", "url", c.url)
		}

		if c.msgChan != nil {
			slog.Info("Closing message channel", "url", c.url)
			close(c.msgChan)
			c.msgChan = nil // prevent double close
			slog.Info("Message channel closed", "url", c.url)
		} else {
			slog.Warn("Message channel already nil, skipping close", "url", c.url)
		}

		c.status = WebSocketDisconnected
		slog.Info("WebSocket disconnected", "url", c.url)
	})
}

func (c *WebSocketApiClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	slog.Debug(fmt.Sprintf("Websocket is currently: %s", c.status))
	return c.status == WebSocketConnected && c.conn != nil
}

func (c *WebSocketApiClient) Send(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("websocket not connected")
	}
	if !c.IsConnected() {
		return fmt.Errorf("websocket is not in connected state")
	}
	err := c.conn.WriteJSON(data)
	if err != nil {
		slog.Error("Error writing JSON to WebSocket", "error", err, "url", c.url)
		return err
	}
	return nil
}

func (c *WebSocketApiClient) Subscribe() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("websocket not connected")
	}

	go func() {
		defer c.Disconnect()

		for {
			select {
			case <-c.ctx.Done():
				slog.Info("Exiting Subscribe loop due to context cancellation",
					"url", c.url)
				return
			case <-c.done: // check for done signal
				slog.Info("Exiting Subscribe loop due to done signal", "url", c.url)
				return
			default:
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(
						err,
						websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
						websocket.CloseAbnormalClosure,
					) {
						slog.Info("Websocket connection closed", "error", err,
							"url", c.url)
					} else {
						slog.Error("Error reading websocket message", "error", err,
							"url", c.url)
					}
					return
				}

				var msg WebSocketMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					slog.Error("Error unmarshaling message", "error", err, "url", c.url)
					continue
				}

				// Send message to channel
				select {
				case c.msgChan <- msg:
				default:
					slog.Warn("Message channel full, dropping message", "url", c.url)
				}

				// Call handler if registered
				if handler, ok := c.handlers[msg.Type]; ok {
					handler(msg.Data)
				}
			}
		}
	}()

	slog.Info("Subscribed to WebSocket", "url", c.url)
	return nil
}

func (c *WebSocketApiClient) OnMessage(msgType string, handler func([]byte)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[msgType] = handler
}

func (c *WebSocketApiClient) MessageChannel() <-chan WebSocketMessage {
	return c.msgChan
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

// WebSocketUpdateMsg represents a message received from the websocket
type WebSocketUpdateMsg struct {
	Type string
	Data interface{}
}

// WebSocketErrorMsg represents an error during websocket communication
type WebSocketErrorMsg struct {
	Error string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	ctx        context.Context
	cancel     context.CancelFunc
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	stateMu    sync.RWMutex
	state      atomic.Value // Store raw JSON
}

func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		ctx:        ctx,
		cancel:     cancel,
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *websocket.Conn, 256),
		unregister: make(chan *websocket.Conn, 256),
		state:      atomic.Value{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			return
		case client := <-h.register:
			h.clients[client] = true
			// Send initial state to new client
			h.stateMu.RLock()
			initialState, _ := h.GetState() // Get raw JSON state
			h.stateMu.RUnlock()
			if initialState != nil {
				message := map[string]interface{}{
					"type": "state",
					"data": initialState,
				}
				if err := client.WriteJSON(message); err != nil {
					slog.Error("Failed to send initial state", "error", err)
					client.Close()
					delete(h.clients, client)
				}
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					slog.Error("Failed to write message", "error", err)
					client.Close()
					delete(h.clients, client)
				}
			}
		}
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "error", err)
		return
	}

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Set ping handler
	conn.SetPingHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
	})

	hub.register <- conn

	go func() {
		defer func() {
			hub.unregister <- conn
		}()

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("Failed to read message", "error", err)
				}
				break
			}

			// Process the message
			var msg map[string]interface{}
			if err := json.Unmarshal(p, &msg); err != nil {
				slog.Error("Failed to unmarshal message", "error", err)
				continue
			}

			msgType, ok := msg["type"].(string)
			if !ok {
				slog.Error("Message type not found or not a string")
				continue
			}

			// Handle different types of messages
			switch msgType {
			case "state":
				// Update state and broadcast to all clients
				data, ok := msg["data"].(interface{})
				if !ok {
					slog.Error("Data not found or invalid")
					continue
				}

				jsonData, err := json.Marshal(data)
				if err != nil {
					slog.Error("Failed to marshal state update", "error", err)
					continue
				}

				hub.stateMu.Lock()
				hub.SetState(jsonData) // Store raw JSON
				hub.stateMu.Unlock()
				hub.broadcast <- jsonData

			case "connected":
				// Handle connection confirmation
				slog.Info("Client connected")

			default:
				slog.Info("Received unknown message type", "type", msgType)
			}
		}
	}()
}

// WebSocketClient wraps the API package's websocket client for TUI integration
type WebSocketClient struct {
	client *WebSocketApiClient
}

// NewWebSocketClient creates a new websocket client
func NewWebSocketClient(url string) *WebSocketClient {
	return &WebSocketClient{
		client: NewWebSocketApiClient(url),
	}
}

// Connect initiates the websocket connection
func (w *WebSocketClient) Connect() tea.Cmd {
	return func() tea.Msg {
		// Initialize websocket client
		slog.Info("Connecting Websocket")

		err := w.client.Connect()
		if err != nil {
			return WebSocketErrorMsg{Error: err.Error()}
		}

		// Start subscribing to messages
		if err := w.client.Subscribe(); err != nil {
			w.client.Disconnect() // Ensure disconnection on subscription failure
			return WebSocketErrorMsg{Error: err.Error()}
		}

		// Start listening for messages
		go func() {
			for msg := range w.client.MessageChannel() {
				var data interface{}
				if err := json.Unmarshal(msg.Data, &data); err != nil {
					slog.Error("Failed to unmarshal message", "error", err)
					continue
				}
				// Send message to TUI
				tea.Cmd(func() tea.Msg {
					return WebSocketUpdateMsg{
						Type: msg.Type,
						Data: data,
					}
				})()
			}
		}()

		return WebSocketUpdateMsg{
			Type: "connected",
			Data: map[string]interface{}{
				"type": "connected",
			},
		}
	}
}

// Disconnect closes the websocket connection
func (w *WebSocketClient) Disconnect() {
	if w.client != nil {
		w.client.Disconnect()
	}
}

// Send sends a message through the websocket
func (w *WebSocketClient) Send(data interface{}) error {
	if w.client == nil {
		return fmt.Errorf("websocket not connected")
	}
	return w.client.Send(data)
}

func (h *Hub) GetState() ([]byte, bool) {
	if state, ok := h.state.Load().([]byte); ok {
		return state, true
	}
	return nil, false
}

func (h *Hub) SetState(newState []byte) {
	h.state.Store(newState)
}

func (w *WebSocketClient) IsConnected() bool {
	if w.client != nil {
		return w.client.IsConnected()
	}
	return false
}

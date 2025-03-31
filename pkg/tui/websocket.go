package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"p1/pkg/api"
	"p1/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

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
	broadcast  chan models.Update
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	stateMu    sync.RWMutex
	state      atomic.Value // Store *models.State
}

func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		ctx:        ctx,
		cancel:     cancel,
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan models.Update, 256),
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
			initialState := models.Update{
				Type: "state",
				Data: h.GetState(),
			}
			h.stateMu.RUnlock()
			if err := client.WriteJSON(initialState); err != nil {
				slog.Error("Failed to send initial state", "error", err)
				client.Close()
				delete(h.clients, client)
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}

		case update := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteJSON(update)
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
			var update models.Update
			err := conn.ReadJSON(&update)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("Failed to read message", "error", err)
				}
				break
			}

			// Handle different types of updates
			switch update.Type {
			case "state":
				// Update state and broadcast to all clients
				hub.stateMu.Lock()
				var stateUpdate models.State
				dataBytes, err := json.Marshal(update.Data)
				if err != nil {
					slog.Error("Failed to marshal state update", "error", err)
					hub.stateMu.Unlock()
					break
				}
				if err := json.Unmarshal(dataBytes, &stateUpdate); err != nil {
					slog.Error("Failed to unmarshal state update", "error", err)
					hub.stateMu.Unlock()
					break
				}
				hub.SetState(&stateUpdate)
				hub.stateMu.Unlock()
				hub.broadcast <- update

			case "connected":
				// Handle connection confirmation
				slog.Info("Client connected")

			default:
				slog.Info("Received unknown message type", "type", update.Type)
			}
		}
	}()
}

// WebSocketClient wraps the API package's websocket client for TUI integration
type WebSocketClient struct {
	client *api.WebSocketClient
}

// NewWebSocketClient creates a new websocket client
func NewWebSocketClient(url string) *WebSocketClient {
	return &WebSocketClient{
		client: api.NewWebSocketClient(url),
	}
}

// Connect initiates the websocket connection
func (w *WebSocketClient) Connect() tea.Cmd {
	return func() tea.Msg {
		err := w.client.Connect()
		if err != nil {
			return VisibleError{message: err.Error()}
		}

		// Start subscribing to messages
		if err := w.client.Subscribe(); err != nil {
			return VisibleError{message: err.Error()}
		}

		// Start listening for messages
		go func() {
			for msg := range w.client.MessageChannel() {
				var data interface{}
				if err := json.Unmarshal(msg.Data, &data); err != nil {
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

func (h *Hub) GetState() *models.State {
	if state, ok := h.state.Load().(*models.State); ok {
		return state
	}
	return nil
}

func (h *Hub) SetState(newState *models.State) {
	h.state.Store(newState)
}

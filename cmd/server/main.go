package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync" // Synchronization primitives (e.g., mutex)

	"github.com/gorilla/websocket" // WebSocket library
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Buffer size for reading WebSocket messages
	WriteBufferSize: 1024, // Buffer size for writing WebSocket messages
	CheckOrigin: func(r *http.Request) bool {
		slog.Info("CheckOrigin called", "origin", r.Header.Get("Origin"))
		return true // Allow all origins (not secure for production)
	},
}

// Hub manages WebSocket connections and broadcasts messages to clients
type Hub struct {
	clients    map[*websocket.Conn]bool // Active WebSocket connections
	broadcast  chan []byte              // Channel for broadcasting updates to clients
	register   chan *websocket.Conn     // Channel for registering new clients
	unregister chan *websocket.Conn     // Channel for unregistering clients
	state      []byte                   // Shared application state as raw JSON
	mutex      sync.RWMutex             // Mutex for synchronizing access to the state
}

// Creates a new Hub instance
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		state:      []byte(`{"servers": [], "brokers": [], "projects": []}`), // Initialize with default JSON
	}
}

// Main loop for the Hub to handle client registration, unregistration, and broadcasting
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register: // Handle new client registration
			h.clients[client] = true
			slog.Info("Client registered", "client", client.RemoteAddr())

			// Send the initial state to the new client
			h.mutex.RLock()
			initialState := map[string]interface{}{
				"type": "state",         // Message type
				"data": string(h.state), // Current state
			}

			h.mutex.RUnlock()

			err := client.WriteJSON(initialState) // Send state as JSON
			if err != nil {
				slog.Error("Failed to send initial state", "error", err, "client", client.RemoteAddr())
			} else {
				slog.Info("Initial state sent to client", "client", client.RemoteAddr())
			}

		case client := <-h.unregister: // Handle client unregistration
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // Remove client from the map
				slog.Info("Client unregistered", "client", client.RemoteAddr())
				err := client.Close() // Close the WebSocket connection
				if err != nil {
					slog.Error("Error closing connection", "error", err, "client", client.RemoteAddr())
				}
			}

		case message := <-h.broadcast: // Broadcast updates to all clients
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message) // Send update as TextMessage
				if err != nil {
					slog.Error("Failed to write message", "error", err, "client", client.RemoteAddr())
					err := client.Close() // Close the connection on error
					if err != nil {
						slog.Error("Error closing connection", "error", err, "client", client.RemoteAddr())
					}
					delete(h.clients, client) // Remove client from the map
				}
			}
		}
	}
}

// Handles WebSocket connections from clients
func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	slog.Info("Attempting to upgrade connection", "remote_addr", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "error", err, "remote_addr", r.RemoteAddr)
		return
	}

	slog.Info("Connection upgraded", "remote_addr", conn.RemoteAddr())

	// Register the new WebSocket connection with the Hub
	hub.register <- conn

	// Goroutine to handle incoming messages from the client
	go func() {
		defer func() {
			hub.unregister <- conn // Unregister the client when the connection closes
			slog.Info("Unregistering client", "remote_addr", conn.RemoteAddr())
		}()

		for {
			_, message, err := conn.ReadMessage() // Read raw message
			if err != nil {
				// Handle unexpected WebSocket closure
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("Failed to read message", "error", err, "remote_addr", conn.RemoteAddr())
				}
				break
			}

			// Process the message
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				slog.Error("Failed to unmarshal message", "error", err, "remote_addr", conn.RemoteAddr())
				continue
			}

			msgType, ok := msg["type"].(string)
			if !ok {
				slog.Error("Message type not found or not a string", "remote_addr", conn.RemoteAddr())
				continue
			}

			// Handle different types of updates from the client
			switch msgType {
			case "state": // Update the shared state and broadcast it
				data, ok := msg["data"].(interface{})
				if !ok {
					slog.Error("Data not found or invalid", "remote_addr", conn.RemoteAddr())
					continue
				}

				jsonData, err := json.Marshal(data)
				if err != nil {
					slog.Error("Failed to marshal state update", "error", err, "remote_addr", conn.RemoteAddr())
					continue
				}

				hub.mutex.Lock()
				hub.state = jsonData // Store raw JSON
				hub.mutex.Unlock()
				hub.broadcast <- jsonData // Broadcast the updated state

			case "connected": // Handle connection confirmation
				slog.Info("Client connected", "remote_addr", conn.RemoteAddr())

			default: // Handle unknown message types
				slog.Info("Received unknown message type", "type", msgType, "remote_addr", conn.RemoteAddr())
			}
		}
	}()
}

// Entry point of the application
func main() {
	hub := newHub() // Create a new Hub instance
	go hub.run()    // Start the Hub's main loop in a separate goroutine

	// Set up the WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	// Start the HTTP server
	addr := ":8080"
	slog.Info("Starting server", "address", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}

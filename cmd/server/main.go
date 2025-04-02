package main

import (
	"log/slog"
	"net/http"
	"p1/pkg/models"
	"sync" // Synchronization primitives (e.g., mutex)

	"github.com/gorilla/websocket" // WebSocket library
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Buffer size for reading WebSocket messages
	WriteBufferSize: 1024, // Buffer size for writing WebSocket messages
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (not secure for production)
	},
}

// Hub manages WebSocket connections and broadcasts messages to clients
type Hub struct {
	clients    map[*websocket.Conn]bool // Active WebSocket connections
	broadcast  chan models.Update       // Channel for broadcasting updates to clients
	register   chan *websocket.Conn     // Channel for registering new clients
	unregister chan *websocket.Conn     // Channel for unregistering clients
	state      *models.State            // Shared application state
	mutex      sync.RWMutex             // Mutex for synchronizing access to the state
}

// Creates a new Hub instance
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan models.Update),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		state: &models.State{
			Servers:  []models.Server{},
			Brokers:  []models.Broker{},
			Projects: []models.Project{},
			Actor:    models.NewActor(),
		},
	}
}

// Main loop for the Hub to handle client registration, unregistration, and broadcasting
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register: // Handle new client registration
			h.clients[client] = true

			// Send the initial state to the new client
			h.mutex.RLock()
			initialState := models.Update{
				Type: "state", // Message type
				Data: h.state, // Current state
			}
			h.mutex.RUnlock()
			client.WriteJSON(initialState) // Send state as JSON

		case client := <-h.unregister: // Handle client unregistration
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // Remove client from the map
				client.Close()            // Close the WebSocket connection
			}

		case update := <-h.broadcast: // Broadcast updates to all clients
			for client := range h.clients {
				err := client.WriteJSON(update) // Send update as JSON
				if err != nil {
					slog.Error("Failed to write message", "error", err)
					client.Close()            // Close the connection on error
					delete(h.clients, client) // Remove client from the map
				}
			}
		}
	}
}

// Handles WebSocket connections from clients
func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "error", err)
		return
	}

	// Register the new WebSocket connection with the Hub
	hub.register <- conn

	// Goroutine to handle incoming messages from the client
	go func() {
		defer func() {
			hub.unregister <- conn // Unregister the client when the connection closes
		}()

		for {
			var update models.Update
			// Read a JSON message from the client
			err := conn.ReadJSON(&update)
			if err != nil {
				// Handle unexpected WebSocket closure
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("Failed to read message", "error", err)
				}
				break
			}

			// Handle different types of updates from the client
			switch update.Type {
			case "state": // Update the shared state and broadcast it
				hub.mutex.Lock()
				if stateData, ok := update.Data.(map[string]interface{}); ok {
					// Update the list of servers in the state
					if servers, ok := stateData["servers"].([]interface{}); ok {
						hub.state.Servers = make([]models.Server, len(servers))
						for i, s := range servers {
							if serverMap, ok := s.(map[string]interface{}); ok {
								hub.state.Servers[i] = models.Server{
									ID:   serverMap["id"].(string),
									Name: serverMap["name"].(string),
									URL:  serverMap["url"].(string),
								}
							}
						}
					}
					// Handle other state updates similarly
				}
				hub.mutex.Unlock()
				hub.broadcast <- update // Broadcast the updated state

			case "connected": // Handle connection confirmation
				slog.Info("Client connected")

			default: // Handle unknown message types
				slog.Info("Received unknown message type", "type", update.Type)
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
	slog.Info("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}

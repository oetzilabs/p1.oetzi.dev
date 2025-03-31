package main

import (
	"log/slog"
	"net/http"
	"sync"

	"p1/pkg/models"

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
	clients    map[*websocket.Conn]bool
	broadcast  chan models.Update
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	state      *models.State
	mutex      sync.RWMutex
}

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

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			// Send initial state to new client
			h.mutex.RLock()
			initialState := models.Update{
				Type: "state",
				Data: h.state,
			}
			h.mutex.RUnlock()
			client.WriteJSON(initialState)

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

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "error", err)
		return
	}

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
				hub.mutex.Lock()
				if stateData, ok := update.Data.(map[string]interface{}); ok {
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

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	slog.Info("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}

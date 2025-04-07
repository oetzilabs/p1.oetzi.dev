package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Service struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Endpoint    string            `json:"endpoint"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

type ServerMetrics struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Storage float64 `json:"storage"`
	Network float64 `json:"network"`
}

type Server struct {
	services   map[string]*Service
	wsUpgrader *websocket.Upgrader
	Address    string
	mu         sync.RWMutex
	srv        *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
	clients    []*string
}

type ServerOptions struct {
	Port string
}

type MessageType string

const (
	TypeListServices    MessageType = "LIST_SERVICES"
	TypeRegisterService MessageType = "REGISTER_SERVICE"
	TypeRemoveService   MessageType = "REMOVE_SERVICE"
	TypeListBrokers     MessageType = "LIST_BROKERS"
	TypeRegisterBroker  MessageType = "REGISTER_BROKER"
	TypeRemoveBroker    MessageType = "REMOVE_BROKER"
	TypeMetrics         MessageType = "METRICS"
	TypeBroadcast       MessageType = "BROADCAST"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
	Sender  string      `json:"sender"`
}

func findOpenPort() string {
	for port := 28080; port <= 38080; port++ {
		if isPortOpen(strconv.Itoa(port)) {
			return strconv.Itoa(port)
		}
	}
	return "28080" // fallback to default if no ports are available
}

func isPortOpen(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func New(options ServerOptions) *Server {
	var port string
	if options.Port == "0" || isPortOpen(port) {
		port = findOpenPort()
	} else {
		port = options.Port
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		services: make(map[string]*Service),
		wsUpgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Address: fmt.Sprintf("localhost:%s", port),
		ctx:     ctx,
		cancel:  cancel,
		clients: []*string{},
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed: " + err.Error())
		return
	}
	defer conn.Close()

	// add the client-id to the list of clients, so we can track them.
	// the client-id is being send with the websocket-handshake
	clientId := r.Header.Get("X-Client-ID")
	if clientId != "" {
		s.clients = append(s.clients, &clientId)
	} else {
		slog.Error("Client-ID not found in websocket-handshake", "error", errors.New("Client-ID not found"))
		return
	}

	// Send periodic health updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			metrics := ServerMetrics{
				CPU:     50.0,
				Memory:  60.0,
				Storage: 70.0,
				Network: 80.0,
			}
			msg := Message{
				Type:    TypeMetrics,
				Payload: metrics,
			}
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			slog.Error("read error:" + err.Error())
			return
		}

		switch msg.Type {
		case TypeListServices:
			s.mu.RLock()
			services := make([]*Service, 0, len(s.services))
			for _, svc := range s.services {
				services = append(services, svc)
			}
			s.mu.RUnlock()

			response := Message{
				Type:    TypeListServices,
				Payload: services,
			}
			conn.WriteJSON(response)

		case TypeRegisterService:
			if service, ok := msg.Payload.(map[string]interface{}); ok {
				// Convert map to Service struct
				svc := &Service{
					ID:          service["id"].(string),
					Name:        service["name"].(string),
					Endpoint:    service["endpoint"].(string),
					Description: service["description"].(string),
					Metadata:    make(map[string]string),
				}
				if metadata, ok := service["metadata"].(map[string]interface{}); ok {
					for k, v := range metadata {
						svc.Metadata[k] = v.(string)
					}
				}

				s.mu.Lock()
				s.services[svc.ID] = svc
				s.mu.Unlock()
			}

		case TypeRemoveService:
			if id, ok := msg.Payload.(string); ok {
				s.mu.Lock()
				delete(s.services, id)
				s.mu.Unlock()
			}
		case TypeBroadcast:
			if msg.Payload == nil {
				slog.Error("msg.Payload is nil")
				return
			}
			payload := msg.Payload.(string)
			for _, clientId := range s.clients {
				if clientId != nil {
					// broadcast to all clients except the sender
					if *clientId != msg.Sender {
						msg := Message{
							Type:    TypeBroadcast,
							Payload: payload,
							Sender:  *clientId,
						}
						// send the message to the client
						if err := conn.WriteJSON(msg); err != nil {
							return
						}
					}
				}
			}
		}
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)

	s.srv = &http.Server{
		Addr:    s.Address,
		Handler: mux,
	}

	go func() {
		<-s.ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
	}()

	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	slog.Info("Running server", "address", s.Address)
	return nil
}

func (s *Server) Shutdown() {
	slog.Info("Shutting down server")
	s.cancel()
}

package client

import (
	"context"
	"p1/pkg/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn        *websocket.Conn
	state       *models.State
	subscribers []chan models.Update
}

func NewClient(url string) *Client {
	// Initialize client
	return &Client{
		state: &models.State{
			Servers:  []models.Server{},
			Brokers:  []models.Broker{},
			Projects: []models.Project{},
			Actor:    models.Actor{},
		},
		conn: nil,
	}
}

func (c *Client) Subscribe() chan models.Update {
	ch := make(chan models.Update)
	c.subscribers = append(c.subscribers, ch)
	return ch
}

func (c *Client) Unsubscribe(ch chan models.Update) {
	for i, subscriber := range c.subscribers {
		if subscriber == ch {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

// Start runs the client in the background
func (c *Client) Start(ctx context.Context) {
	// Create WebSocket connection
	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(ctx, "ws://localhost:8080/ws", nil)
	if err != nil {
		// TODO: Handle connection error
		return
	}
	c.conn = conn

	// Start message handling goroutine
	go func() {
		defer conn.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var update models.Update
				err := conn.ReadJSON(&update)
				if err != nil {
					// Handle read error
					return
				}

				// Update state based on the received update
				switch update.Type {
				case "state":
					if state, ok := update.Data.(models.State); ok {
						c.state = &state
					}
				case "server":
					if server, ok := update.Data.(models.Server); ok {
						// Update or add server
						found := false
						for i, s := range c.state.Servers {
							if s.ID == server.ID {
								c.state.Servers[i] = server
								found = true
								break
							}
						}
						if !found {
							c.state.Servers = append(c.state.Servers, server)
						}
					}
				case "broker":
					if broker, ok := update.Data.(models.Broker); ok {
						// Update or add broker
						found := false
						for i, b := range c.state.Brokers {
							if b.ID == broker.ID {
								c.state.Brokers[i] = broker
								found = true
								break
							}
						}
						if !found {
							c.state.Brokers = append(c.state.Brokers, broker)
						}
					}
				case "project":
					if project, ok := update.Data.(models.Project); ok {
						// Update or add project
						found := false
						for i, p := range c.state.Projects {
							if p.ID == project.ID {
								c.state.Projects[i] = project
								found = true
								break
							}
						}
						if !found {
							c.state.Projects = append(c.state.Projects, project)
						}
					}
				case "actor":
					if actor, ok := update.Data.(models.Actor); ok {
						c.state.Actor = actor
					}
				}

				// Notify all subscribers of the update
				for _, ch := range c.subscribers {
					select {
					case ch <- update:
					default:
						// Skip if subscriber is not ready to receive
					}
				}
			}
		}
	}()
}

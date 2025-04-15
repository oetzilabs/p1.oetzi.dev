package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"p1/pkg/messages"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Endpoint    string            `json:"endpoint"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

type Memory struct {
	MemTotal uint64 `json:"mem_total"` // Total memory in bytes
	MemFree  uint64 `json:"mem_free"`  // Free memory in bytes
}

type Network struct {
	ID           string `json:"id"`            // Interface identifier
	RxBytes      uint64 `json:"rx_bytes"`      // Received bytes
	RxPackets    uint64 `json:"rx_packets"`    // Received packets
	RxErrors     uint64 `json:"rx_errors"`     // Received errors
	RxDropped    uint64 `json:"rx_dropped"`    // Received dropped packets
	RxFifo       uint64 `json:"rx_fifo"`       // Received FIFO errors
	RxFrame      uint64 `json:"rx_frame"`      // Received frame errors
	RxCompressed uint64 `json:"rx_compressed"` // Received compressed packets
	RxMulticast  uint64 `json:"rx_multicast"`  // Received multicast packets
	TxBytes      uint64 `json:"tx_bytes"`      // Transmitted bytes
	TxPackets    uint64 `json:"tx_packets"`    // Transmitted packets
	TxErrors     uint64 `json:"tx_errors"`     // Transmitted errors
	TxDropped    uint64 `json:"tx_dropped"`    // Transmitted dropped packets
	TxFifo       uint64 `json:"tx_fifo"`       // Transmitted FIFO errors
	TxFrame      uint64 `json:"tx_frame"`      // Transmitted frame errors
	TxCompressed uint64 `json:"tx_compressed"` // Transmitted compressed packets
	TxMulticast  uint64 `json:"tx_multicast"`  // Transmitted multicast packets
}

type Storage struct {
	Disks []Disk `json:"disks"` // List of disks
}

type Disk struct {
	MountPoint string `json:"mount_point"` // Mount point of the disk
	Total      uint64 `json:"total"`       // Total space in bytes
	Used       uint64 `json:"used"`        // Used space in bytes
}

type CPU struct {
	User   uint64 `json:"user"`   // User CPU time
	System uint64 `json:"system"` // System CPU time
	Idle   uint64 `json:"idle"`   // Idle CPU time
	Total  uint64 `json:"total"`  // Total CPU time
}

type ServerMetrics struct {
	CPU     *CPU     `json:"cpu"`     // CPU metrics
	Memory  *Memory  `json:"memory"`  // Memory metrics
	Storage *Storage `json:"storage"` // Storage metrics
	Network *Network `json:"network"` // Network metrics
}

type Server struct {
	ID         string
	services   map[string]*Service // Map of registered services
	wsUpgrader *websocket.Upgrader // WebSocket upgrader
	Address    string              // Server address
	WSLink     string              // WebSocket link
	mu         sync.RWMutex        // Mutex for protecting service map
	srv        *http.Server        // HTTP server
	ctx        context.Context     // Context for server lifecycle management
	cancel     context.CancelFunc  // Cancel function for context
	clients    []*string           // List of connected client IDs
}

type ServerOptions struct {
	Port string // Port number to listen on
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
		ID:       uuid.New().String(),
		services: make(map[string]*Service),
		wsUpgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Address: fmt.Sprintf("localhost:%s", port),
		WSLink:  fmt.Sprintf("ws://localhost:%s/ws", port),
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
	clientId := r.Header.Get("CID")

	if clientId != "" {
		s.clients = append(s.clients, &clientId)
	}

	// Send periodic health updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			metrics := ServerMetrics{
				CPU:     getCPU(),
				Memory:  getMemory(),
				Storage: getStorage(),
				Network: getNetwork(),
			}
			msg := messages.Message{
				Type:    messages.TypeMetrics,
				Payload: metrics,
				Sender:  s.ID,
			}
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}()

	for {
		var msg messages.Message
		if err := conn.ReadJSON(&msg); err != nil {
			slog.Error("read error:" + err.Error())
			return
		}

		switch msg.Type {
		case messages.TypeListServices:
			// List registered services and send them to the client.
			s.mu.RLock()
			services := make([]*Service, 0, len(s.services))
			for _, svc := range s.services {
				services = append(services, svc)
			}
			s.mu.RUnlock()

			response := messages.Message{
				Type:    messages.TypeListServices,
				Payload: services,
			}
			conn.WriteJSON(response)

		case messages.TypeRegisterService:
			// Register a new service.
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

		case messages.TypeRemoveService:
			// Remove a service.
			if id, ok := msg.Payload.(string); ok {
				s.mu.Lock()
				delete(s.services, id)
				s.mu.Unlock()
			}
		case messages.TypeBroadcast:
			// Broadcast a message to all connected clients.
			if msg.Payload == nil {
				slog.Error("msg.Payload is nil")
				return
			}
			payload := msg.Payload.(string)
			for _, clientId := range s.clients {
				if clientId != nil {
					// broadcast to all clients except the sender
					if *clientId != msg.Sender {
						msg := messages.Message{
							Type:    messages.TypeBroadcast,
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
	time.Sleep(500 * time.Millisecond)
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown() {
	slog.Info("Shutting down server")
	s.cancel()
}

// getCPU retrieves CPU usage statistics from /proc/stat.
func getCPU() *CPU {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		slog.Error("Failed to read /proc/stat", "error", err)
		return nil
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] != "cpu" {
			continue
		}

		user, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			slog.Error("Failed to parse user CPU time", "error", err)
			return nil
		}

		system, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			slog.Error("Failed to parse system CPU time", "error", err)
			return nil
		}

		idle, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			slog.Error("Failed to parse idle CPU time", "error", err)
			return nil
		}

		total := user + system + idle
		return &CPU{
			User:   user,
			System: system,
			Idle:   idle,
			Total:  total,
		}
	}

	slog.Error("No CPU stats found")
	return nil
}

// getMemory retrieves memory usage statistics from /proc/meminfo.
func getMemory() *Memory {
	contents, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		slog.Error("Failed to read /proc/meminfo", "error", err)
		return nil
	}

	lines := strings.Split(string(contents), "\n")
	memInfo := make(map[string]uint64)

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		// Remove " kB" and parse the value
		valueStr = strings.ReplaceAll(valueStr, " kB", "")
		value, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			slog.Error("Failed to parse memory value", "error", err, "key", key)
			continue
		}

		memInfo[key] = value
	}

	memTotal, ok1 := memInfo["MemTotal"]
	memFree, ok2 := memInfo["MemFree"]

	if !ok1 || !ok2 {
		slog.Error("Missing MemTotal or MemFree in /proc/meminfo")
		return nil
	}

	return &Memory{
		MemTotal: memTotal,
		MemFree:  memFree,
	}
}

// getStorage retrieves storage usage statistics using the "df" command.
func getStorage() *Storage {
	cmd := exec.Command("df", "-BG") // Use "df -BG" to get sizes in GB
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Error executing df command", "error", err)
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var disks []Disk

	// Skip header line
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) < 6 {
			continue
		}

		mountPoint := fields[5]
		totalStr := strings.ReplaceAll(fields[1], "G", "") // Remove "G" suffix
		usedStr := strings.ReplaceAll(fields[2], "G", "")  // Remove "G" suffix

		total, err := strconv.ParseUint(totalStr, 10, 64)
		if err != nil {
			slog.Error("Failed to parse total storage", "error", err, "value", totalStr)
			continue
		}

		used, err := strconv.ParseUint(usedStr, 10, 64)
		if err != nil {
			slog.Error("Failed to parse used storage", "error", err, "value", usedStr)
			continue
		}

		disks = append(disks, Disk{
			MountPoint: mountPoint,
			Total:      total * 1024 * 1024 * 1024, // Convert GB to Bytes
			Used:       used * 1024 * 1024 * 1024,  // Convert GB to Bytes
		})
	}

	return &Storage{
		Disks: disks,
	}
}

// getNetwork retrieves network interface statistics from /proc/net/dev.
func getNetwork() *Network {
	contents, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		slog.Error("Failed to read /proc/net/dev", "error", err)
		return nil
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 17 {
			continue
		}
		if !strings.Contains(fields[0], "lo") {
			continue
		}
		//remove ":" from the interface name
		networkInterface := strings.ReplaceAll(fields[0], ":", "")
		rxBytes, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxbytes", "error", err)
			return nil
		}

		rxPackets, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxpackets", "error", err)
			return nil
		}

		rxErrors, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxerrors", "error", err)
			return nil
		}

		rxDropped, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxdropped", "error", err)
			return nil
		}

		rxFifo, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxfifo", "error", err)
			return nil
		}

		rxFrame, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxframe", "error", err)
			return nil
		}

		rxCompressed, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxcompressed", "error", err)
			return nil
		}

		rxMulticast, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			slog.Error("Failed to parse rxmulticast", "error", err)
			return nil
		}

		txBytes, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txbytes", "error", err)
			return nil
		}

		txPackets, err := strconv.ParseUint(fields[10], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txpackets", "error", err)
			return nil
		}

		txErrors, err := strconv.ParseUint(fields[11], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txerrors", "error", err)
			return nil
		}

		txDropped, err := strconv.ParseUint(fields[12], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txdropped", "error", err)
			return nil
		}

		txFifo, err := strconv.ParseUint(fields[13], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txfifo", "error", err)
			return nil
		}

		txFrame, err := strconv.ParseUint(fields[14], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txframe", "error", err)
			return nil
		}

		txCompressed, err := strconv.ParseUint(fields[15], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txcompressed", "error", err)
			return nil
		}

		txMulticast, err := strconv.ParseUint(fields[16], 10, 64)
		if err != nil {
			slog.Error("Failed to parse txmulticast", "error", err)
			return nil
		}

		return &Network{
			ID:           networkInterface,
			RxBytes:      rxBytes,
			RxPackets:    rxPackets,
			RxErrors:     rxErrors,
			RxDropped:    rxDropped,
			RxFifo:       rxFifo,
			RxFrame:      rxFrame,
			RxCompressed: rxCompressed,
			RxMulticast:  rxMulticast,
			TxBytes:      txBytes,
			TxPackets:    txPackets,
			TxErrors:     txErrors,
			TxDropped:    txDropped,
			TxFifo:       txFifo,
			TxFrame:      txFrame,
			TxCompressed: txCompressed,
			TxMulticast:  txMulticast,
		}
	}

	return nil
}

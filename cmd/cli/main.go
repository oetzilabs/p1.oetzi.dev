package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"p1/pkg/api"
	"p1/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	log, err := os.Create("output.log")
	if err != nil {
		panic(err)
	}
	defer log.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(log, &slog.HandlerOptions{})))

	// Start websocket server in a goroutine
	go func() {
		hub := api.NewHub()
		go hub.Run()

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			api.HandleWebSocket(hub, w, r)
		})

		slog.Info("Starting websocket server on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			slog.Error("Failed to start websocket server", "error", err)
		}
	}()

	model, err := tui.NewModel(lipgloss.DefaultRenderer(), "ws://localhost:8080/ws", []string{})
	if err != nil {
		panic(err)
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

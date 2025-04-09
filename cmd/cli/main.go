package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"p1/pkg/client"
	"p1/pkg/config"
	"p1/pkg/server"
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

	config := config.New()

	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)
	tuiDone := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var srv *server.Server
	var cl *client.Client
	if config.WithServer {
		wg.Add(1)
		serverOptions := server.ServerOptions{
			Port: config.ServerPort,
		}
		srv = server.New(serverOptions)

		go func() {
			defer wg.Done()
			if err := srv.Start(); err != nil {
				slog.Error("Error starting server", "error", err)
				os.Exit(1)
			}
		}()
	}

	// Unified shutdown handler
	go func() {
		select {
		case <-sigChan:
			slog.Info("Received shutdown signal")
		case <-tuiDone:
			slog.Info("TUI closed")
		}
		if srv != nil {
			srv.Shutdown()
		}
		if cl != nil {
			err := cl.Stop()
			if err != nil {
				slog.Error("Error stopping client", "error", err.Error())
			}
		}
		// Ensure clean exit
		if config.WithTui {
			os.Exit(0)
		}
	}()

	if config.WithTui {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(tuiDone)

			cl = client.NewClient(srv.WSLink)
			err := cl.Init()
			if err != nil {
				slog.Error("Error initializing client", "error", err.Error())
				os.Exit(1)
			}
			err = cl.Start()
			if err != nil {
				slog.Error("Error starting client", "error", err.Error())
				os.Exit(1)
			}

			model := tui.NewModel(lipgloss.DefaultRenderer(), cl)
			if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
				slog.Error("Error running TUI", "error", err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
}

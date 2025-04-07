package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	WithTui    bool
	WithServer bool
	ServerPort string
}

const ENV_TUI = "TUI"
const ENV_SERVER = "SERVER"
const ENV_PORT = "PORT"

const FLAG_NO_TUI = "no-tui"
const FLAG_NO_SERVER = "no-server"
const FLAG_PORT = "port"

func New() *Config {
	cfg := &Config{
		WithTui:    true,
		WithServer: true,
		ServerPort: "0",
	}

	// Environment variables take precedence over defaults
	if v := os.Getenv(ENV_TUI); v != "" {
		cfg.WithTui = parseBool(v)
	}
	if v := os.Getenv(ENV_SERVER); v != "" {
		cfg.WithServer = parseBool(v)
	}
	if v := os.Getenv(ENV_PORT); v != "" {
		cfg.ServerPort = v
	}

	// Command line flags take precedence over environment variables
	flag.BoolVar(&cfg.WithTui, FLAG_NO_TUI, !cfg.WithTui, "disable TUI")
	flag.BoolVar(&cfg.WithServer, FLAG_NO_SERVER, !cfg.WithServer, "disable server")
	flag.StringVar(&cfg.ServerPort, FLAG_PORT, cfg.ServerPort, "server port")
	flag.Parse()

	// Invert the "no-" flags
	cfg.WithTui = !cfg.WithTui
	cfg.WithServer = !cfg.WithServer

	return cfg
}

func parseBool(v string) bool {
	b, _ := strconv.ParseBool(v)
	return b
}

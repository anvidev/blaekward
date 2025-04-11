package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anvidev/blaekward/internal/conn"
	"github.com/anvidev/blaekward/internal/database"
	"github.com/anvidev/blaekward/internal/env"
	"github.com/anvidev/blaekward/internal/message"
)

type Server struct {
	config config
	logger *slog.Logger

	cmu     sync.RWMutex
	clients map[string]*Client

	pmu     sync.Mutex
	pending map[string]chan message.ResponseData
}

type Client struct {
	conn     *conn.WrappedConn
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Printers []message.Printer `json:"printers"`
}

type config struct {
	socketPort string
	httpPort   string
	database   databaseConfig
}

type databaseConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func main() {
	cfg := config{
		socketPort: env.GetString("SOCKET_PORT", ":9091"),
		httpPort:   env.GetString("HTTP_PORT", ":9090"),
		database: databaseConfig{
			addr:         env.GetString("POSTGRES_ADDR", "postgres://admin:adminpassword@localhost/blaekward?sslmode=disable"),
			maxOpenConns: env.GetInt("POSTGRES_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("POSTGRES_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetDuration("POSTGRES_MAX_IDLE_TIME", time.Minute*15),
		},
	}

	s := &Server{
		config:  cfg,
		logger:  slog.Default(),
		clients: make(map[string]*Client),
		pending: make(map[string]chan message.ResponseData),
	}

	_, err := database.NewPostgres(
		cfg.database.addr,
		cfg.database.maxOpenConns,
		cfg.database.maxIdleConns,
		cfg.database.maxIdleTime,
	)
	if err != nil {
		s.logger.Error("database connection failed", "error", err)
		return
	}

	go s.startApi(cfg.httpPort)
	go s.startSocket(cfg.socketPort)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	s.logger.Info("shutting down servers")
}

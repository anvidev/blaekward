package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anvidev/blaekward/internal/database"
	"github.com/anvidev/blaekward/internal/env"
	"github.com/anvidev/blaekward/internal/event"
	"github.com/anvidev/blaekward/internal/store"
)

func main() {
	cfg := config{
		httpPort: env.String("HTTP_PORT", ":9090"),
		database: databaseConfig{
			addr:         env.String("POSTGRES_ADDR", "postgres://admin:adminpassword@localhost/db?sslmode=disable"),
			maxOpenConns: env.Int("POSTGRES_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.Int("POSTGRES_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.Duration("POSTGRES_MAX_IDLE_TIME", time.Minute*15),
		},
	}

	logger := slog.Default()
	db, err := database.NewPostgres(
		cfg.database.addr,
		cfg.database.maxOpenConns,
		cfg.database.maxIdleConns,
		cfg.database.maxIdleTime,
	)

	if err != nil {
		logger.Error("database connection failed", "error", err)
		return
	}

	store := store.New(db)

	s := &Server{
		config:  cfg,
		logger:  logger,
		store:   store,
		clients: make(map[int64]*Client),
		pending: make(map[string]chan event.Event),
	}

	go s.startApi()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	s.logger.Info("shutting down servers")
}

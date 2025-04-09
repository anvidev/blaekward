package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anvidev/blaekward/internal/message"
)

type Server struct {
	logger *slog.Logger

	cmu     sync.RWMutex
	clients map[string]*Client

	pmu     sync.Mutex
	pending map[string]chan message.ResponseData
}

type Client struct {
	conn     net.Conn
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Printers []message.Printer `json:"printers"`
}

func main() {
	s := &Server{
		logger:  slog.Default(),
		clients: make(map[string]*Client),
		pending: make(map[string]chan message.ResponseData),
	}

	go s.startApi(":9090")
	go s.startSocket(":9091")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	slog.Info("shutting down servers")
}

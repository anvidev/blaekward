package main

import (
	"net"
	"sync"

	"go.uber.org/zap"
)

var (
	socketPort = ":9091"
	apiPort    = ":9090"
)

type server struct {
	logger *zap.SugaredLogger

	wg sync.WaitGroup

	mu      sync.RWMutex
	clients map[string]*client
}

type client struct {
	conn net.Conn

	msgCh chan string
}

func new() *server {
	return &server{
		clients: make(map[string]*client),
	}
}

func (s *server) start() {
	s.wg.Add(2)
	go s.startSocket()
	go s.startApi()
	s.wg.Wait()
}

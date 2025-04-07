package main

import (
	"net"
	"sync"

	"github.com/anvidev/blaekward/internal/event"
	"github.com/anvidev/blaekward/internal/utils"
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

	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Printers []event.Printer `json:"printers"`

	msgCh chan utils.Response
}

// type printer struct {
// 	Name string `json:"name"`
// 	Addr string `json:"addr"`
// }

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

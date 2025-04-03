package main

import (
	"bufio"
	"net"

	"github.com/google/uuid"
)

func (s *server) startSocket() {
	ln, err := net.Listen("tcp", socketPort)
	if err != nil {
		s.logger.Fatalw("failed to start socket", "error", err)
	}

	defer func() {
		ln.Close()
		s.wg.Done()
	}()

	s.logger.Infow("started socket", "port", socketPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Errorw("error accepting connection", "error", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	id := uuid.New().String()

	cl := &client{conn: conn, msgCh: make(chan string, 10)}

	s.mu.Lock()
	s.clients[id] = cl
	s.mu.Unlock()

	s.logger.Infow("client connected", "id", id)

	defer func() {
		s.mu.Lock()
		conn.Close()
		delete(s.clients, id)
		s.mu.Unlock()
		s.logger.Infow("client disconnected", "id", id)
	}()

	scanner := bufio.NewScanner(cl.conn)

	for scanner.Scan() {
		b := scanner.Bytes()

		cl.msgCh <- string(b)

		s.logger.Infow("message from client", "message", string(b))
	}
}

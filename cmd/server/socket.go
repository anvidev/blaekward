package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/anvidev/blaekward/internal/event"
	"github.com/anvidev/blaekward/internal/utils"
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
	id := uuid.New().String()

	cl := &client{
		ID:    id,
		conn:  conn,
		msgCh: make(chan utils.Response, 10),
	}

	if err := s.authenticate(conn, cl); err != nil {
		s.logger.Errorw("error authenticating client", "error", err)
		return
	}

	s.mu.Lock()
	s.clients[id] = cl
	s.mu.Unlock()

	s.logger.Infow("client connected", "id", id)

	defer func() {
		s.mu.Lock()
		conn.Close()
		delete(s.clients, id)
		s.mu.Unlock()
		s.wg.Done()
		s.logger.Infow("client disconnected", "id", id)
	}()

	scanner := bufio.NewScanner(cl.conn)

	for scanner.Scan() {
		b := scanner.Bytes()

		var resp utils.Response

		if err := json.Unmarshal(b, &resp); err != nil {
			s.logger.Infow("failed to unmarshal respoonse", "error", err)
			resp.Success = false
			resp.Message = "failed to unmarshal response"
			cl.msgCh <- resp
			continue
		}

		cl.msgCh <- resp

		s.logger.Infow("message from client", "message", resp.Message)
	}
}

func (s *server) authenticate(conn net.Conn, client *client) error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	var authEvent event.AuthenticationEvent

	if err := authEvent.Decode(buf[:n]); err != nil {
		return err
	}

	if authEvent.Key != "abc123" {
		return fmt.Errorf("wrong authentication")
	}

	client.Printers = authEvent.Printers
	client.Name = authEvent.Name

	return nil
}

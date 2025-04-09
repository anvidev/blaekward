package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/anvidev/blaekward/internal/conn"
	"github.com/anvidev/blaekward/internal/message"
	"github.com/google/uuid"
)

func (s *Server) startSocket(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("failed to listen on socket", "error", err)
		return
	}

	s.logger.Info("socket started", "port", addr)

	for {
		rawConn, err := ln.Accept()
		if err != nil {
			s.logger.Error("accept error", "error", err)
			continue
		}

		go s.handleConnection(conn.Wrap(rawConn))
	}
}

func (s *Server) handleConnection(c *conn.WrappedConn) {
	defer c.Close()

	clientID := uuid.New().String()
	client := &Client{ID: clientID, conn: c}

	if err := s.authenticate(client); err != nil {
		s.logger.Error("failed to authenticate", "error", err)
		return
	}

	s.cmu.Lock()
	s.clients[clientID] = client
	s.cmu.Unlock()

	s.logger.Info("client connected", "id", clientID)

	scanner := bufio.NewScanner(client.conn)

	for scanner.Scan() {
		data := scanner.Bytes()
		// data, err := c.ReadMessage()
		// if err != nil {
		// 	switch err {
		// 	case io.EOF:
		// 		s.logger.Info("client disconnected")
		// 		return
		// 	default:
		// 		s.logger.Error("read error", "error", err)
		// 		continue
		// 	}
		// }

		msg, err := message.Decode(data)
		if err != nil {
			s.logger.Error("decode error", "error", err)
			continue
		}

		switch msg.Type {
		// case message.TypeAuth:
		// 	s.handleAuth(msg.Data, client)
		case message.TypeResponse:
			s.handleResponse(msg.Data)
		default:
			s.logger.Warn("unhandled message type", "type", msg.Type)
		}
	}

	if err := scanner.Err(); err != nil {
		s.logger.Error("connection error", "error", err)
	}
}

func (s *Server) authenticate(c *Client) error {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return err
	}

	msg, err := message.Decode(buf[:n])
	if err != nil {
		s.logger.Error("failed to decode auth message", "error", err)
		return err
	}

	var data message.AuthData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		s.logger.Error("failed to decode auth", "error", err)
		return err
	}

	if data.Key != "abc123" {
		return fmt.Errorf("wrong authentication")
	}

	c.Printers = data.Printers
	c.Name = data.Name

	return nil
}

func (s *Server) handleResponse(raw json.RawMessage) {
	var data message.ResponseData
	if err := json.Unmarshal(raw, &data); err != nil {
		s.logger.Error("failed to decode response", "error", err)
		return
	}

	s.pmu.Lock()
	ch, found := s.pending[data.RequestID]
	s.pmu.Unlock()

	if !found {
		s.logger.Error("request channel not found", "request_id", data.RequestID)
		return
	}

	ch <- data
}

func (s *Server) handleAuth(raw json.RawMessage, client *Client) {
	var data message.AuthData
	if err := json.Unmarshal(raw, &data); err != nil {
		s.logger.Error("failed to decode auth", "error", err)
		return
	}

	if data.Key != "abc123" {
		s.logger.Warn("invalid auth attempt", "client", data.Name)
		return
	}

	s.cmu.Lock()
	client.Printers = data.Printers
	s.cmu.Unlock()

	s.logger.Info("client authenticated", "client", data.Name)
}

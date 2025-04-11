package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/anvidev/blaekward/internal/conn"
	"github.com/anvidev/blaekward/internal/id"
	"github.com/anvidev/blaekward/internal/message"
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
	clientID := id.Generate(16, id.LettersLower, id.LetterUpper, id.Numbers)
	client := &Client{ID: clientID, conn: c}

	if err := s.authenticate(client); err != nil {
		s.logger.Error("failed to authenticate", "error", err)
		return
	}

	s.cmu.Lock()
	s.clients[clientID] = client
	s.cmu.Unlock()

	s.logger.Info("client connected", "id", clientID)

	defer func() {
		s.cmu.Lock()
		c.Close()
		delete(s.clients, clientID)
		s.cmu.Unlock()
		s.logger.Info("client removed from server")
	}()

	for {
		data, err := c.ReadMessage()
		if err != nil {
			switch err {
			case io.ErrUnexpectedEOF:
			case io.EOF:
				s.logger.Info("client disconnected")
				return
			default:
				s.logger.Error("read error", "error", err)
				continue
			}
		}

		msg, err := message.Decode(data)
		if err != nil {
			s.logger.Error("decode error", "error", err)
			continue
		}

		switch msg.Type {
		case message.TypeResponse:
			s.handleResponse(msg)
		default:
			s.logger.Warn("unhandled message type", "type", msg.Type)
		}
	}
}

func (s *Server) authenticate(c *Client) error {
	c.conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	defer c.conn.SetReadDeadline(time.Time{})

	bytes, err := c.conn.ReadMessage()
	if err != nil {
		s.logger.Error("failed to decode auth message", "error", err)
		return err
	}

	msg, err := message.Decode(bytes)
	if err != nil {
		s.logger.Error("failed to decode auth message", "error", err)
		return err
	}

	if msg.Type != message.TypeAuth {
		s.logger.Error("failed to decode auth message", "error", "wrong message type")
		return fmt.Errorf("failed to decode auth message: %s", msg.Type)
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

func (s *Server) handleResponse(msg *message.Message) {
	var data message.ResponseData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		s.logger.Error("failed to decode response", "error", err)
		return
	}

	s.pmu.Lock()
	ch, found := s.pending[msg.RequestID]
	s.pmu.Unlock()

	if !found {
		s.logger.Error("request channel not found", "request_id", msg.RequestID)
		return
	}

	ch <- data
}

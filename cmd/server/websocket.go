package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anvidev/blaekward/internal/event"
	ws "github.com/gorilla/websocket"
)

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}
	defer conn.Close()

	s.logger.Info("incoming connection", "remote", conn.RemoteAddr())

	cid, auth, err := s.authenticate(conn)
	if err != nil {
		s.logger.Error("authentication failed", "error", err)
		return
	}

	s.cmu.Lock()
	s.clients[cid] = &Client{
		conn:     conn,
		ID:       cid,
		Name:     auth.Name,
		Printers: auth.Printers,
	}
	s.cmu.Unlock()

	defer func(s *Server, conn *ws.Conn, cid int64) {
		s.cmu.Lock()
		delete(s.clients, cid)
		s.cmu.Unlock()
		s.logger.Info("client disconnected", "client_id", cid)
	}(s, conn, cid)

	s.logger.Info("client authenticated", "client_id", cid)

	if err := s.sendSuccess(conn, "authenticated"); err != nil {
		return
	}

	s.read(conn, cid)
}

func (s *Server) authenticate(conn *ws.Conn) (int64, *event.AuthData, error) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetReadDeadline(time.Time{})

	var evt event.Event
	if err := conn.ReadJSON(&evt); err != nil {
		s.sendError(conn, event.CodeBadRequest, err.Error(), true)
		return -1, nil, err
	}

	if evt.Type != event.TypeAuth {
		msg := fmt.Sprintf("incorrect event type received for authentication: %s", evt.Type)
		s.sendError(conn, event.CodeBadRequest, msg, true)
		return -1, nil, fmt.Errorf(msg)
	}

	var auth event.AuthData
	if err := json.Unmarshal(evt.Data, &auth); err != nil {
		return -1, nil, err
	}

	ctx := context.Background()

	apikey, err := s.store.APIKeys.GetByKey(ctx, auth.Key)
	if err != nil {
		return -1, nil, err
	}

	return apikey.CustomerID, &auth, nil
}

func (s *Server) read(conn *ws.Conn, clientID int64) {
	for {
		var evt event.Event
		if err := conn.ReadJSON(&evt); err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				s.logger.Error("unexpected close", "client_id", clientID, "error", err)
			}
			return
		}

		s.pmu.Lock()
		respCh, found := s.pending[evt.RequestID]
		s.pmu.Unlock()

		if !found {
			s.logger.Warn("orphaned request", "request_id", evt.RequestID, "event_type", evt.Type)
			continue
		}

		respCh <- evt
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/anvidev/blaekward/internal/contextkeys"
	"github.com/anvidev/blaekward/internal/event"
	"github.com/anvidev/blaekward/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type Server struct {
	config config
	logger *slog.Logger
	store  store.Store

	cmu     sync.RWMutex
	clients map[int64]*Client

	pmu     sync.Mutex
	pending map[string]chan event.Event
}

type Client struct {
	conn     *ws.Conn
	ID       int64
	Name     string
	Printers []event.Printer
}

type config struct {
	httpPort string
	database databaseConfig
}

type databaseConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

var upgrader = ws.Upgrader{}

func (s *Server) startApi() {
	mux := chi.NewMux()

	srv := &http.Server{
		Addr:    s.config.httpPort,
		Handler: mux,
	}

	s.mount(mux)

	s.logger.Info("api started", "port", s.config.httpPort)

	if err := srv.ListenAndServe(); err != nil {
		s.logger.Error("failed to listen on api", "error", err)
	}
}

func (s *Server) mount(r *chi.Mux) {
	r.Use(allowAllCORS)
	r.Use(requestID)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", handlePing)
		r.Get("/ws", s.handleWebsocket)

		r.Route("/admin", func(r chi.Router) {
			r.Post("/api-keys", s.handleCreateAPIKey)
		})

		r.Group(func(r chi.Router) {
			r.Use(s.customerAuthorization)
			r.Get("/printers", s.handleListPrinters)
			r.Post("/message", s.handleMessage)
		})
	})

}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	cid, ok := r.Context().Value(contextkeys.CustomerID).(int64)
	if !ok {
		s.internalServerError(w, r, fmt.Errorf("customer id not found in context"))
		return
	}

	var evt event.Event
	if err := readJSON(w, r, &evt); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	reqID := r.Header.Get(requestIDTargetHeader)
	if reqID == "" {
		reqID = uuid.NewString()
	}
	evt.RequestID = reqID

	s.cmu.RLock()
	defer s.cmu.RUnlock()
	client, found := s.clients[cid]
	if !found {
		s.notFoundError(w, r, fmt.Errorf("client with id %d was not found", cid))
		return
	}

	responseCh := make(chan event.Event, 1)

	s.pmu.Lock()
	s.pending[reqID] = responseCh
	s.pmu.Unlock()

	defer func() {
		s.pmu.Lock()
		delete(s.pending, reqID)
		s.pmu.Unlock()
	}()

	if err := client.conn.WriteJSON(&evt); err != nil {
		s.internalServerError(w, r, err)
		return
	}

	select {
	case resp := <-responseCh:
		s.handleEvent(w, r, resp)
	case <-time.After(10 * time.Second):
		s.gatewayTimeout(w, r, fmt.Errorf("no response from client"))
	}
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request, evt event.Event) {
	switch evt.Type {
	case event.TypeError:
		var evtData event.ErrorData
		if err := event.Extract(evt.Data, &evtData); err != nil {
			s.internalServerError(w, r, err)
			return
		}
		if evtData.Code == event.CodeNotFound {
			s.notFoundError(w, r, fmt.Errorf(evtData.Message))
		} else if evtData.Code == event.CodeNotReachable {
			s.notFoundError(w, r, fmt.Errorf(evtData.Message))
		} else {
			s.internalServerError(w, r, fmt.Errorf(evtData.Message))
		}
	case event.TypeSuccess:
		var evtData event.SuccessData
		if err := event.Extract(evt.Data, &evtData); err != nil {
			s.internalServerError(w, r, err)
			return
		}
		type response struct {
			Message string `json:"message"`
		}
		if err := writeJSON(w, http.StatusOK, response{Message: evtData.Message}); err != nil {
			s.internalServerError(w, r, err)
			return
		}
	default:
		s.internalServerError(w, r, fmt.Errorf("unknown event %s received on server for request %s", evt.Type, evt.RequestID))
	}
}

func (s *Server) handleListPrinters(w http.ResponseWriter, r *http.Request) {
	cid, ok := r.Context().Value(contextkeys.CustomerID).(int64)
	if !ok {
		s.internalServerError(w, r, fmt.Errorf("customer id not found in request context"))
		return
	}

	s.cmu.RLock()
	defer s.cmu.RUnlock()

	printers := []event.Printer{}

	for _, client := range s.clients {
		if client.ID == cid {
			printers = append(printers, client.Printers...)
		}
	}

	type response struct {
		Printers []event.Printer `json:"printers"`
	}

	if err := writeJSON(w, http.StatusOK, response{Printers: printers}); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *Server) sendError(conn *ws.Conn, code event.ErrorCode, msg string, fatal bool) error {
	data, err := json.Marshal(event.ErrorData{
		Code:    code,
		Message: msg,
		Fatal:   fatal,
	})
	if err != nil {
		return err
	}
	evt := event.New(event.TypeError, data)
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetWriteDeadline(time.Time{})
	return conn.WriteJSON(evt)
}

func (s *Server) sendSuccess(conn *ws.Conn, msg string) error {
	data, err := json.Marshal(event.SuccessData{
		Message: msg,
	})
	if err != nil {
		return err
	}
	evt := event.New(event.TypeSuccess, data)
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetWriteDeadline(time.Time{})
	return conn.WriteJSON(evt)
}

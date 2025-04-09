package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anvidev/blaekward/internal/message"
	"github.com/google/uuid"
)

func (s *Server) startApi(addr string) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s.mount(mux)

	s.logger.Info("api started", "port", addr)

	if err := srv.ListenAndServe(); err != nil {
		s.logger.Error("failed to listen on api", "error", err)
	}
}

func (s *Server) mount(m *http.ServeMux) {
	m.HandleFunc("GET /ping", handlePing)
	m.HandleFunc("GET /clients", s.handleListClients)
	m.HandleFunc("POST /{id}/message", s.handleMessage)
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var msg message.Message

	if err := readJSON(w, r, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	requestID := uuid.New().String()
	msg.RequestID = requestID

	s.logger.Info("message received", "type", msg.Type, "request_id", msg.RequestID)

	s.cmu.RLock()
	client, found := s.clients[id]
	s.cmu.RUnlock()

	if !found {
		http.Error(w, "client not found", http.StatusNotFound)
		return
	}

	responseCh := make(chan message.ResponseData, 1)

	s.pmu.Lock()
	s.pending[requestID] = responseCh
	s.pmu.Unlock()

	defer func() {
		s.pmu.Lock()
		delete(s.pending, requestID)
		s.pmu.Unlock()
	}()

	bytes, err := message.Encode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := client.conn.Write(append(bytes, '\n')); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	select {
	case resp := <-responseCh:
		if resp.Success {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp.Message)
		} else {
			http.Error(w, resp.Message, http.StatusInternalServerError)
		}
	case <-time.After(10 * time.Second):
		http.Error(w, "print timeout", http.StatusGatewayTimeout)
	}
}

func (s *Server) handleListClients(w http.ResponseWriter, r *http.Request) {
	s.cmu.RLock()
	defer s.cmu.RUnlock()

	clients := []*Client{}

	for _, client := range s.clients {
		clients = append(clients, client)
	}

	type response struct {
		Clients []*Client `json:"clients"`
	}

	if err := writeJSON(w, http.StatusOK, response{Clients: clients}); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

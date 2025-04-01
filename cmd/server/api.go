package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anvidev/inkbyte/internal/event"
)

func (s *server) startApi() {
	defer s.wg.Done()

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         apiPort,
		Handler:      mux,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
		WriteTimeout: time.Second * 30,
	}

	s.mount(mux)

	s.logger.Infow("started api", "port", apiPort)

	if err := srv.ListenAndServe(); err != nil {
		s.logger.Fatalw("failed to start api", "error", err)
	}
}

func (s *server) mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /ping", s.pingHandler)
	mux.HandleFunc("GET /machines", s.listMachinesHandler)
	mux.HandleFunc("POST /{id}/print", s.printHandler)
}

type response struct {
	Machines []string `json:"machines"`
}

func (s *server) listMachinesHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	machines := []string{}

	for id := range s.clients {
		machines = append(machines, id)
	}

	if err := writeJSON(w, http.StatusOK, response{Machines: machines}); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *server) printHandler(w http.ResponseWriter, r *http.Request) {
	var evt event.Event

	if err := readJSON(w, r, &evt); err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown"):
			s.badRequestError(w, r, err)
		default:
			s.internalServerError(w, r, err)
		}
		return
	}

	if err := validateJSON(&evt); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	data, err := json.Marshal(evt)
	if err != nil {
		s.badRequestError(w, r, err)
		return
	}

	data = append(data, '\n')

	id := r.PathValue("id")

	s.mu.RLock()
	defer s.mu.RUnlock()
	cl, found := s.clients[id]
	if !found {
		s.notFoundError(w, r, fmt.Errorf("machine not found"))
		return
	}

	if _, err := cl.conn.Write(data); err != nil {
		s.internalServerError(w, r, err)
		return
	}

	resp := <-cl.msgCh

	w.Write([]byte(resp))
}

func (s *server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

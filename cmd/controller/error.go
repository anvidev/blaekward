package main

import "net/http"

type errorEnvelope struct {
	Error string `json:"error"`
}

func (s *Server) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Error("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusInternalServerError, &errorEnvelope{Error: "internal server error"})
}

func (s *Server) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Info("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusBadRequest, &errorEnvelope{Error: err.Error()})
}

func (s *Server) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warn("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusNotFound, &errorEnvelope{Error: err.Error()})
}

func (s *Server) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Error("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusConflict, &errorEnvelope{Error: err.Error()})
}

func (s *Server) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusUnauthorized, &errorEnvelope{Error: err.Error()})
}

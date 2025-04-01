package main

import "net/http"

type errorEnvelope struct {
	Error string `json:"error"`
}

func (s *server) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusInternalServerError, &errorEnvelope{Error: "internal server error"})
}

func (s *server) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Infow("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusBadRequest, &errorEnvelope{Error: err.Error()})
}

func (s *server) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusNotFound, &errorEnvelope{Error: err.Error()})
}

func (s *server) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusConflict, &errorEnvelope{Error: err.Error()})
}

func (s *server) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warnw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSON(w, http.StatusUnauthorized, &errorEnvelope{Error: err.Error()})
}

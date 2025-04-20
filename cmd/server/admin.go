package main

import (
	"net/http"

	"github.com/anvidev/blaekward/internal/store"
)

type createAPIKeyRequest struct {
	CustomerID int64 `json:"customer_id"`
}

func (s *Server) handleCreateAPIKey(w http.ResponseWriter, r *http.Request) {
	var payload createAPIKeyRequest
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	apikey, err := s.store.APIKeys.Create(ctx, payload.CustomerID)
	if err != nil {
		switch err {
		case store.ErrDublicateCustomerID:
			s.conflictError(w, r, err)
		default:
			s.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, apikey); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

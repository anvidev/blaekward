package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/anvidev/blaekward/internal/contextkeys"
	"github.com/anvidev/blaekward/internal/id"
	"github.com/anvidev/blaekward/internal/store"
)

var (
	requestIDTargetHeader = "X-Request-Identifier"
)

func (s *Server) customerAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			s.unauthorizedError(w, r, fmt.Errorf("authorization header is malformored"))
			return
		}

		authKey := authParts[1]
		ctx := r.Context()

		apikey, err := s.store.APIKeys.GetByKey(ctx, authKey)
		if err != nil {
			switch err {
			case store.ErrAPIKeyNotFound:
				s.unauthorizedError(w, r, err)
			default:
				s.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, contextkeys.CustomerID, apikey.CustomerID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(requestIDTargetHeader)
		if reqID == "" {
			reqID = id.String(52, id.LettersLower, id.LettersUpper, id.Numbers)
			r.Header.Set(requestIDTargetHeader, reqID)
		}
		next.ServeHTTP(w, r)
	})
}

func allowAllCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

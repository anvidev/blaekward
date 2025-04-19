package store

import (
	"context"
	"database/sql"
	"time"
)

const (
	DefaultQueryTimeoutDuration = time.Second * 30
)

type Store struct {
	APIKeys APIKeysStorer
}

func New(db *sql.DB) Store {
	return Store{
		APIKeys: &apiKeyStore{db},
	}
}

type APIKeysStorer interface {
	Create(context.Context, int64) (*APIKey, error)
	GetByKey(context.Context, string) (*APIKey, error)
}

package store

import (
	"context"
	"time"
)

type Store struct {
	Key KeyStorer
}

const (
	DefaultQueryTimeoutDuration = time.Second * 30
)

type KeyStorer interface {
	Create(context.Context, int64) (*Key, error)
	GetByKey(context.Context, string) (*Key, error)
}

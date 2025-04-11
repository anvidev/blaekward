package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/agentstation/uuidkey"
	"github.com/google/uuid"
)

type Key struct {
	ID         int64
	CustomerID int64
	Key        string
	Inserted   time.Time
}

type keyStore struct {
	db *sql.DB
}

func (s *keyStore) Create(ctx context.Context, customerID int64) (*Key, error) {
	baseUUID := uuid.New().String()
	apikey, err := uuidkey.Encode(baseUUID)
	if err != nil {
		return nil, err
	}

	key := &Key{CustomerID: customerID, Key: apikey.String()}

	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeoutDuration)
	defer cancel()

	stmt := `insert into bl_keys (customer_id, key) values ($1, $2) returning id, inserted`

	err = s.db.QueryRowContext(ctx, stmt, key.CustomerID, key.Key).Scan(&key.ID, &key.Inserted)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *keyStore) GetByKey(ctx context.Context, key string) (*Key, error) {
	stmt := `select id, customer_id, key, inserted from bl_keys where key = $1`

	baseKey, err := uuidkey.Parse(key)
	if err != nil {
		return nil, err
	}

	uuidKey, err := baseKey.UUID()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeoutDuration)
	defer cancel()

	var ckey Key

	err = s.db.QueryRowContext(ctx, stmt, baseKey.UUID()).Scan()
}

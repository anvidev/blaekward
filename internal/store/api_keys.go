package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/agentstation/uuidkey"
	"github.com/google/uuid"
)

var (
	ErrAPIKeyNotFound      = errors.New("api key not found")
	ErrDublicateCustomerID = errors.New("customer id already has api key")
)

type APIKey struct {
	CustomerID int64     `json:"customer_id"`
	Key        string    `json:"key"`
	Inserted   time.Time `json:"inserted"`
	Expiry     time.Time `json:"expiry"`
}

type apiKeyStore struct {
	db *sql.DB
}

func (s *apiKeyStore) Create(ctx context.Context, customerID int64, expiryDur time.Duration) (*APIKey, error) {
	uuid := uuid.New().String()
	expiry := time.Now().Add(expiryDur)

	encoded, err := uuidkey.Encode(uuid)
	if err != nil {
		return nil, err
	}

	key := &APIKey{
		CustomerID: customerID,
		Key:        encoded.String(),
		Inserted:   time.Now(),
		Expiry:     expiry,
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeoutDuration)
	defer cancel()

	stmt := `insert into api_keys (customer_id, uuid, expiry) values ($1, $2, $3) returning inserted`

	err = s.db.QueryRowContext(ctx, stmt, key.CustomerID, uuid, expiry).Scan(&key.Inserted)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "api_keys_customer_id_key"`:
			return nil, ErrDublicateCustomerID
		default:
			return nil, err
		}
	}

	return key, nil
}

func (s *apiKeyStore) GetByKey(ctx context.Context, key string) (*APIKey, error) {
	encoded, err := uuidkey.Parse(key)
	if err != nil {
		return nil, err
	}

	uuid, err := encoded.UUID()
	if err != nil {
		return nil, err
	}

	stmt := `select uuid, customer_id, inserted, expiry from api_keys where uuid = $1`

	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeoutDuration)
	defer cancel()

	var APIKey APIKey

	err = s.db.QueryRowContext(
		ctx,
		stmt,
		uuid,
	).Scan(
		&APIKey.Key,
		&APIKey.CustomerID,
		&APIKey.Inserted,
		&APIKey.Expiry,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrAPIKeyNotFound
		default:
			return nil, err
		}
	}

	return &APIKey, nil
}

func (k *APIKey) IsExpired() bool {
	if k.Expiry.IsZero() {
		return false // never expires
	}
	return time.Now().After(k.Expiry)
}

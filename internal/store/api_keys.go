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
}

type apiKeyStore struct {
	db *sql.DB
}

func (s *apiKeyStore) Create(ctx context.Context, customerID int64) (*APIKey, error) {
	uuid := uuid.New().String()

	encoded, err := uuidkey.Encode(uuid)
	if err != nil {
		return nil, err
	}

	key := &APIKey{
		CustomerID: customerID,
		Key:        encoded.String(),
		Inserted:   time.Now().Local(),
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeoutDuration)
	defer cancel()

	stmt := `insert into api_keys (customer_id, uuid) values ($1, $2) returning inserted`

	err = s.db.QueryRowContext(ctx, stmt, key.CustomerID, uuid).Scan(&key.Inserted)
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

	stmt := `select uuid, customer_id, inserted from api_keys where uuid = $1`

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

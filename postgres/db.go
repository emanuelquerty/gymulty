package postgres

import (
	"context"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.Store = (*Store)(nil)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func Connect(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

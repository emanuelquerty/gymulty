package postgres

import (
	"context"
	"fmt"

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

func Connect(dbname string, dbusername string, dbpassword string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbusername, dbpassword, dbname)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

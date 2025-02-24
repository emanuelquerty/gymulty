package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
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
	err := createDBIfNotExists(dbname, dbusername, dbpassword)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbusername, dbpassword, dbname)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func createDBIfNotExists(dbname string, dbusername string, dbpassword string) error {
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/postgres", dbusername, dbpassword)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	var exists bool
	query := "SELECT EXISTS (SELECT FROM pg_database WHERE datname=$1)"
	err = conn.QueryRow(ctx, query, dbname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if !exists {
		query := fmt.Sprintf("CREATE DATABASE %s", dbname)
		if _, err := conn.Exec(ctx, query); err != nil {
			return fmt.Errorf("error creating database: %w", err)
		} else {
			slog.Info("Database created successfully")
		}
	}
	return nil
}

package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/emanuelquerty/gymulty/config"
	"github.com/emanuelquerty/gymulty/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var _ domain.Store = (*Store)(nil)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func Connect(conf config.DBconfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", conf.DBusername, conf.DBpassword, conf.DBname)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func CreateDBIfNotExists(conf config.DBconfig) error {
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/postgres", conf.DBusername, conf.DBpassword)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	var exists bool
	query := "SELECT EXISTS (SELECT FROM pg_database WHERE datname=$1)"
	err = conn.QueryRow(ctx, query, conf.DBname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if !exists {
		query := fmt.Sprintf("CREATE DATABASE %s", conf.DBname)
		if _, err := conn.Exec(ctx, query); err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	}
	return nil
}

func RunMigrations(embedMigrations embed.FS, conf config.DBconfig) error {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", conf.DBusername, conf.DBpassword, conf.DBname)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("error opening database for migrations: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting goose dialect: %w", err)
	}

	if err := goose.Up(db, "postgres/migrations"); err != nil {
		return fmt.Errorf("error running up migration: %w", err)
	}
	return nil
}

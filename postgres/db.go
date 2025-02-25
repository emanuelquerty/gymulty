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

func Connect(conf config.Database) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func CreateDBIfNotExists(conf config.Database) error {
	ctx := context.Background()
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=postgres",
		conf.User, conf.Password, conf.Host, conf.Port)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	var exists bool
	query := "SELECT EXISTS (SELECT FROM pg_database WHERE datname=$1)"
	err = conn.QueryRow(ctx, query, conf.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if !exists {
		query := fmt.Sprintf("CREATE DATABASE %s", conf.Name)
		if _, err := conn.Exec(ctx, query); err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	}
	return nil
}

func RunMigrations(embedMigrations embed.FS, conf config.Database) error {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)

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

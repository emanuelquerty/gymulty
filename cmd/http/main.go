package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/emanuelquerty/gymulty"
	"github.com/emanuelquerty/gymulty/http"
	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/emanuelquerty/gymulty/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	conf := loadConfigFromEnv(logger)
	dbpool, err := postgres.Connect(conf.dbname, conf.dbusername, conf.dbpassword)
	if err != nil {
		log.Fatal(err)
	}

	err = runMigrations(gymulty.EmbedMigrations, *conf)
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServer(dbpool, logger)

	server.Use(middleware.Logger)
	server.Use(middleware.SetHeader("Content-Type", "application/json"))
	server.Use(middleware.AddRequestID)

	err = server.ListenAndServe(8080)
	if err != nil {
		log.Fatal(err)
	}
}

func runMigrations(embedMigrations embed.FS, conf config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", conf.dbusername, conf.dbpassword, conf.dbname)
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

type config struct {
	dbname     string
	dbusername string
	dbpassword string
}

func loadConfigFromEnv(logger *slog.Logger) *config {
	conf := new(config)
	conf.dbname = getEnv(logger, "DB_NAME")
	conf.dbusername = getEnv(logger, "DB_USERNAME")
	conf.dbpassword = getEnv(logger, "DB_PASSWORD")
	return conf
}

func getEnv(logger *slog.Logger, name string) string {
	env, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("environment variable: could not find %s", name)
		logger.Error(err.Error())
	}
	return env
}

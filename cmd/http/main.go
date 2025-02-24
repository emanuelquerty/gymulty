package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/emanuelquerty/gymulty"
	"github.com/emanuelquerty/gymulty/config"
	"github.com/emanuelquerty/gymulty/http"
	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/emanuelquerty/gymulty/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbconfig := config.LoadDB(logger)

	err = postgres.CreateDBIfNotExists(*dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Database created successfully!")
	err = postgres.RunMigrations(gymulty.EmbedMigrations, *dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	dbpool, err := postgres.Connect(*dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Database connected successfully!")

	server := http.NewServer(dbpool, logger)

	server.Use(middleware.Logger)
	server.Use(middleware.SetHeader("Content-Type", "application/json"))
	server.Use(middleware.AddRequestID)

	err = server.ListenAndServe(8080)
	if err != nil {
		log.Fatal(err)
	}
}

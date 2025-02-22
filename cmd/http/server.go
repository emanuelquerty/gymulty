package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/emanuelquerty/gymulty/http"
	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/emanuelquerty/gymulty/postgres"
)

func main() {
	dsn := "postgres://postgres:lealdade@localhost:5432/gymulty"
	dbpool, err := postgres.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := http.NewServer(dbpool, logger)

	server.Use(middleware.Logger)
	server.Use(middleware.AddRequestID)

	err = server.ListenAndServe(8080)
	if err != nil {
		log.Fatal(err)
	}
}

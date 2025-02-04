package main

import (
	"log/slog"
	"os"

	"github.com/emanuelquerty/gymulty/http"
	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/emanuelquerty/gymulty/postgres"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := "postgres://postgres:lealdade@localhost:5432/gymulty"
	conn, err := postgres.Connect(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server := http.NewServer(conn, logger)

	server.Use(middleware.Logger)

	if err = server.ListenAndServe("8080"); err != nil {
		logger.Error(err.Error())
	}
}

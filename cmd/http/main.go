package main

import (
	"log/slog"
	"os"

	"github.com/emanuelquerty/multency/http"
	"github.com/emanuelquerty/multency/postgres"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := "postgres://postgres:lealdade@localhost:5432/multenc"
	conn, err := postgres.Connect(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server := http.NewServer(conn, logger)
	if err = server.ListenAndServe("8080"); err != nil {
		logger.Error(err.Error())
	}
}

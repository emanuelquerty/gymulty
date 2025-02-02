package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/emanuelquerty/multency/http"
	"github.com/emanuelquerty/multency/postgres"
)

func main() {
	dsn := "postgres://postgres:lealdade@localhost:5432/multency"
	conn, err := postgres.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	server := http.NewServer(conn, logger)
	err = server.ListenAndServe("8080")
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/emanuelquerty/multency/http"
	"github.com/emanuelquerty/multency/postgres"
)

func main() {
	dsn := "postgres://postgres:lealdade@localhost:5432/multency"
	conn, err := postgres.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServer(conn)
	err = server.ListenAndServe("8080")
	if err != nil {
		log.Fatal(err)
	}
}

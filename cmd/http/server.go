package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/emanuelquerty/gymulty/http"
	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/emanuelquerty/gymulty/postgres"
	"github.com/joho/godotenv"
)

type Config struct {
	dbname     string
	dbusername string
	dbpassword string
}

func newConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	conf := new(Config)
	conf.dbname = getEnv("DB_NAME")
	conf.dbusername = getEnv("DB_USERNAME")
	conf.dbpassword = getEnv("DB_PASSWORD")
	return conf
}

func getEnv(name string) string {
	env, ok := os.LookupEnv(name)
	if !ok {
		log.Fatal("Environment variable: could not find database name")
	}
	return env
}

func main() {
	conf := newConfig()
	dbpool, err := postgres.Connect(conf.dbname, conf.dbusername, conf.dbpassword)
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

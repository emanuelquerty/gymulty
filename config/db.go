package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadDB(logger *slog.Logger) *Database {
	conf := new(Database)

	conf.Host = getEnv(logger, "DB_HOST")
	conf.Port = getEnv(logger, "DB_PORT")
	conf.Name = getEnv(logger, "DB_NAME")
	conf.User = getEnv(logger, "DB_USER")
	conf.Password = getEnv(logger, "DB_PASSWORD")
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

package config

import (
	"fmt"
	"log/slog"
	"os"
)

type DBconfig struct {
	DBname     string
	DBusername string
	DBpassword string
}

func LoadDB(logger *slog.Logger) *DBconfig {
	conf := new(DBconfig)
	conf.DBname = getEnv(logger, "DB_NAME")
	conf.DBusername = getEnv(logger, "DB_USERNAME")
	conf.DBpassword = getEnv(logger, "DB_PASSWORD")
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

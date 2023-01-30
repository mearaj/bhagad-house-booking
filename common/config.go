package common

import (
	"os"
)

type Config struct {
	DatabaseDriver string
	DatabaseURL    string
}

// These values map to the environment variable names. Should be same as above config file
const (
	DatabaseDriver = "DB_DRIVER"
	DatabaseURL    = "DATABASE_URL"
)

const DefaultDatabaseDriver = "mongodb"
const DefaultDatabaseURL = "mongodb://root:example@localhost:27017/?"

func LoadConfig() (config Config) {
	if config.DatabaseDriver = os.Getenv(DatabaseDriver); config.DatabaseDriver == "" {
		config.DatabaseDriver = DefaultDatabaseDriver
	}
	if config.DatabaseURL = os.Getenv(DatabaseURL); config.DatabaseURL == "" {
		config.DatabaseURL = DefaultDatabaseURL
	}
	return config
}

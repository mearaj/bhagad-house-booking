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

const DefaultDatabaseDriver = "postgres"
const DefaultDatabaseURL = "postgresql://root:secret@localhost:5432/bhagad_house_booking?sslmode=disable"

func LoadConfig() (config Config) {
	if config.DatabaseDriver = os.Getenv(DatabaseDriver); config.DatabaseDriver == "" {
		config.DatabaseDriver = DefaultDatabaseDriver
	}
	if config.DatabaseURL = os.Getenv(DatabaseURL); config.DatabaseURL == "" {
		config.DatabaseURL = DefaultDatabaseURL
	}
	return config
}

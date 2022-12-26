// Package config defines the structs needed for the server, database and other future components that need configurations
package config

import (
	"github.com/rs/zerolog/log"
	"os"
)

type ServerConf struct {
	Port string
}
type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// GetEnv returns default value if the env variable is not found.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Warn().Msgf("error occurred while trying to read %s env variable, it will be set to default value %s", key, fallback)
	return fallback
}

// Package config defines the structs needed for the server, database and other future components that need configurations
package config

import (
	"github.com/rs/zerolog/log"
	"os"
)

var Config Configuration

type Configuration struct {
	Server ServerConfig
	DB     DbConfig
	Auth   AuthConfig
}

type ServerConfig struct {
	Port string
}
type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type AuthConfig struct {
	Username string
	Password string
}

func BuildConfig() {
	conf := Configuration{
		Server: ServerConfig{Port: GetEnv("PORT", "8080")},
		DB: DbConfig{
			Host:     GetEnv("POSTGRES_HOST", "127.0.0.1"),
			User:     GetEnv("POSTGRES_USER", "user"),
			Password: GetEnv("POSTGRES_PASSWORD", "password"),
			Name:     GetEnv("POSTGRES_DB", "tasks"),
			Port:     GetEnv("POSTGRES_PORT", "5432"),
		},
		Auth: AuthConfig{
			Username: os.Getenv("APP_USERNAME"),
			Password: os.Getenv("APP_PASSWORD"),
		},
	}
	Config = conf
}

// GetEnv returns default value if the env variable is not found.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Warn().Msgf("error occurred while trying to read %s env variable, it will be set to default value %s", key, fallback)
	return fallback
}

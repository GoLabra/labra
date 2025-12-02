package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	// MinJWTSecretLength is the minimum required length for JWT secrets
	// HS256 requires at least 32 bytes (256 bits) for security, but we use 16 as a practical minimum
	MinJWTSecretLength = 16
)

type Config struct {
	DSN                  string `env:"DSN,required"`
	DBDialect            string `env:"DB_DIALECT,required"`
	ServerPort           string `env:"SERVER_PORT"`
	EntSchemaPath        string `env:"ENT_SCHEMA_PATH"`
	SecretKey            string `env:"SECRET_KEY,required"`
	CentrifugoApiAddress string `env:"CENTRIFUGO_API_ADDRESS,required"`
	CentrifugoKey        string `env:"CENTRIFUGO_API_KEY,required"`
	FileStorageProvider  string `env:"FILE_STORAGE_PROVIDER"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
}

func New() (*Config, error) {
	godotenv.Load()

	var cfg = &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse environment variables: %w", err)
	}

	// Set defaults for optional configuration if not provided
	if cfg.ServerPort == "" {
		cfg.ServerPort = "4000"
	}
	if cfg.EntSchemaPath == "" {
		cfg.EntSchemaPath = "./ent/schema"
	}
	if cfg.FileStorageProvider == "" {
		cfg.FileStorageProvider = "local"
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = "./storage"
	}

	// Validate JWT secret length
	if len(cfg.SecretKey) < MinJWTSecretLength {
		return nil, fmt.Errorf("SECRET_KEY must be at least %d characters long for security", MinJWTSecretLength)
	}

	return cfg, nil
}

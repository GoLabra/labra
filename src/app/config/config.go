package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DSN                  string `env:"DSN,required"`
	DBDialect            string `env:"DB_DIALECT,required"`
	ServerPort           string `env:"SERVER_PORT,required"`
	EntSchemaPath        string `env:"ENT_SCHEMA_PATH,required"`
	SecretKey            string `env:"SECRET_KEY,required"`
	CentrifugoApiAddress string `env:"CENTRIFUGO_API_ADDRESS,required"`
	CentrifugoKey        string `env:"CENTRIFUGO_API_KEY,required"`
}

func New() (*Config, error) {
	godotenv.Load()

	var cfg = &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse environment variables: %w", err)
	}

	return cfg, nil
}

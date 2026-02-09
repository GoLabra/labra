package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	// MinJWTSecretLength is the minimum required length for JWT secrets
	// HS256 requires at least 32 bytes (256 bits) for security, but we use 16 as a practical minimum
	MinJWTSecretLength = 16
)

// Config holds non-sensitive application configuration loaded from environment variables.
// Sensitive credentials (DSN, SecretKey, CentrifugoKey) are managed separately via the secrets package.
type Config struct {
	DBDialect            string `env:"DB_DIALECT,required"`
	ServerPort           string `env:"SERVER_PORT"`
	EntSchemaPath        string `env:"ENT_SCHEMA_PATH"`
	CentrifugoApiAddress string `env:"CENTRIFUGO_API_ADDRESS,required"`
	FileStorageProvider  string `env:"FILE_STORAGE_PROVIDER"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	// Environment specifies which Infisical environment to fetch secrets from (dev, staging, prod)
	Environment        string `env:"APP_ENVIRONMENT"`
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS"`
	// CSPPolicy allows overriding the Content-Security-Policy header per deployment.
	CSPPolicy string `env:"CSP_POLICY"`
}

// Secrets holds sensitive credentials fetched from Infisical or environment variables.
// These values should never be logged or written to disk.
type Secrets struct {
	DSN           string
	SecretKey     string
	CentrifugoKey string
}

// AppConfig combines non-sensitive configuration with sensitive secrets.
// This is the main configuration struct used by the application.
type AppConfig struct {
	Config
	Secrets
}

// New loads non-sensitive configuration from environment variables.
// Use NewWithSecrets to get a complete AppConfig with secrets loaded.
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
	if cfg.Environment == "" {
		cfg.Environment = "dev"
	}
	if cfg.CORSAllowedOrigins == "" {
		cfg.CORSAllowedOrigins = "http://localhost:3000"
	}
	if cfg.CSPPolicy == "" {
		cfg.CSPPolicy = "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'"
	}

	return cfg, nil
}

// ValidateSecrets validates that secrets meet security requirements.
func ValidateSecrets(s *Secrets) error {
	if s.DSN == "" {
		return fmt.Errorf("DSN is required")
	}
	if s.SecretKey == "" {
		return fmt.Errorf("SECRET_KEY is required")
	}
	if len(s.SecretKey) < MinJWTSecretLength {
		return fmt.Errorf("SECRET_KEY must be at least %d characters long for security", MinJWTSecretLength)
	}
	if s.CentrifugoKey == "" {
		return fmt.Errorf("CENTRIFUGO_API_KEY is required")
	}
	return nil
}

// NewAppConfig creates a complete application configuration by combining
// non-sensitive config with the provided secrets.
func NewAppConfig(cfg *Config, secrets *Secrets) (*AppConfig, error) {
	if err := ValidateSecrets(secrets); err != nil {
		return nil, err
	}

	return &AppConfig{
		Config:  *cfg,
		Secrets: *secrets,
	}, nil
}

// CORSAllowedOriginsList parses CORS_ALLOWED_ORIGINS (comma-separated) into a list.
func (c *Config) CORSAllowedOriginsList() []string {
	raw := strings.TrimSpace(c.CORSAllowedOrigins)
	if raw == "" {
		return []string{"http://localhost:3000"}
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		o := strings.TrimSpace(p)
		if o != "" {
			out = append(out, o)
		}
	}
	return out
}

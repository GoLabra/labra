// Package secrets provides an abstraction layer for fetching sensitive credentials
// from various secret management backends (Infisical, environment variables, etc.).
//
// This package follows a provider pattern, allowing the application to switch between
// different secret backends without changing application code.
package secrets

import (
	"context"
	"fmt"

	"github.com/GoLabra/labra/config"
)

// Provider defines the interface for fetching secrets from a backend.
// Implementations should handle authentication and caching as appropriate.
type Provider interface {
	// GetSecrets fetches all required secrets for the specified environment.
	// The environment parameter maps to Infisical environments (dev, staging, prod).
	GetSecrets(ctx context.Context, environment string) (*config.Secrets, error)

	// Name returns the provider name for logging purposes.
	Name() string
}

// ProviderConfig holds configuration for secret providers.
type ProviderConfig struct {
	// Infisical configuration
	InfisicalClientID     string
	InfisicalClientSecret string
	InfisicalProjectID    string
	InfisicalSiteURL      string // Optional: defaults to https://app.infisical.com

	// Environment fallback (for testing/legacy support)
	UseEnvFallback bool
}

// NewProviderFromEnv creates a new secret provider based on environment configuration.
// It returns an Infisical provider if credentials are available, otherwise falls back to env.
func NewProviderFromEnv() (Provider, error) {
	cfg := loadProviderConfig()

	// If Infisical credentials are available, use Infisical
	if cfg.InfisicalClientID != "" && cfg.InfisicalClientSecret != "" {
		return NewInfisicalProvider(cfg)
	}

	// Fall back to environment variables
	return NewEnvProvider(), nil
}

// LoadSecrets is a convenience function that loads secrets using the appropriate provider.
// It automatically selects the provider based on environment configuration.
func LoadSecrets(ctx context.Context, environment string) (*config.Secrets, error) {
	provider, err := NewProviderFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to create secret provider: %w", err)
	}

	secrets, err := provider.GetSecrets(ctx, environment)
	if err != nil {
		return nil, fmt.Errorf("failed to load secrets via %s: %w", provider.Name(), err)
	}

	return secrets, nil
}


package secrets

import (
	"context"
	"fmt"
	"os"

	"github.com/GoLabra/labra/config"
)

// EnvProvider fetches secrets from environment variables.
// This is used as a fallback for local development, testing, or legacy deployments.
//
// WARNING: This provider reads secrets directly from environment variables.
// For production use, prefer InfisicalProvider for better security and auditability.
type EnvProvider struct{}

// NewEnvProvider creates a new environment variable secret provider.
func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

// Name returns the provider name for logging purposes.
func (p *EnvProvider) Name() string {
	return "env"
}

// GetSecrets fetches secrets from environment variables.
// The environment parameter is ignored since env vars don't support multi-environment.
func (p *EnvProvider) GetSecrets(ctx context.Context, environment string) (*config.Secrets, error) {
	secrets := &config.Secrets{}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DSN environment variable is required")
	}
	secrets.DSN = dsn

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY environment variable is required")
	}
	secrets.SecretKey = secretKey

	centrifugoKey := os.Getenv("CENTRIFUGO_API_KEY")
	if centrifugoKey == "" {
		return nil, fmt.Errorf("CENTRIFUGO_API_KEY environment variable is required")
	}
	secrets.CentrifugoKey = centrifugoKey

	return secrets, nil
}

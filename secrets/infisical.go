package secrets

import (
	"context"
	"fmt"
	"os"

	"github.com/GoLabra/labra/config"
	infisical "github.com/infisical/go-sdk"
)

// InfisicalProvider fetches secrets from Infisical using Machine Identity authentication.
// It implements the Provider interface for secret management.
type InfisicalProvider struct {
	client    infisical.InfisicalClientInterface
	projectID string
	siteURL   string
}

// NewInfisicalProvider creates a new Infisical secret provider.
// It requires valid Machine Identity credentials (client ID and secret).
func NewInfisicalProvider(cfg ProviderConfig) (*InfisicalProvider, error) {
	if cfg.InfisicalClientID == "" {
		return nil, fmt.Errorf("INFISICAL_CLIENT_ID is required for Infisical provider")
	}
	if cfg.InfisicalClientSecret == "" {
		return nil, fmt.Errorf("INFISICAL_CLIENT_SECRET is required for Infisical provider")
	}
	if cfg.InfisicalProjectID == "" {
		return nil, fmt.Errorf("INFISICAL_PROJECT_ID is required for Infisical provider")
	}

	siteURL := cfg.InfisicalSiteURL
	if siteURL == "" {
		siteURL = "https://app.infisical.com"
	}

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl:          siteURL,
		AutoTokenRefresh: true,
	})

	// Authenticate using Universal Auth (Machine Identity)
	_, err := client.Auth().UniversalAuthLogin(cfg.InfisicalClientID, cfg.InfisicalClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Infisical: %w", err)
	}

	return &InfisicalProvider{
		client:    client,
		projectID: cfg.InfisicalProjectID,
		siteURL:   siteURL,
	}, nil
}

// Name returns the provider name for logging purposes.
func (p *InfisicalProvider) Name() string {
	return "infisical"
}

// GetSecrets fetches all required secrets from Infisical for the specified environment.
// The environment parameter should match your Infisical environment slug (dev, staging, prod).
func (p *InfisicalProvider) GetSecrets(ctx context.Context, environment string) (*config.Secrets, error) {
	// Fetch all secrets from the root path
	secretsList, err := p.client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:          p.projectID,
		Environment:        environment,
		SecretPath:         "/",
		AttachToProcessEnv: false, // We don't want to pollute the process env
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list secrets from Infisical: %w", err)
	}

	// Build a map for easy lookup
	secretsMap := make(map[string]string)
	for _, secret := range secretsList {
		secretsMap[secret.SecretKey] = secret.SecretValue
	}

	// Extract required secrets
	secrets := &config.Secrets{}

	dsn, ok := secretsMap["DSN"]
	if !ok {
		return nil, fmt.Errorf("DSN secret not found in Infisical environment '%s'", environment)
	}
	secrets.DSN = dsn

	secretKey, ok := secretsMap["SECRET_KEY"]
	if !ok {
		return nil, fmt.Errorf("SECRET_KEY secret not found in Infisical environment '%s'", environment)
	}
	secrets.SecretKey = secretKey

	centrifugoKey, ok := secretsMap["CENTRIFUGO_API_KEY"]
	if !ok {
		return nil, fmt.Errorf("CENTRIFUGO_API_KEY secret not found in Infisical environment '%s'", environment)
	}
	secrets.CentrifugoKey = centrifugoKey

	return secrets, nil
}

// loadProviderConfig loads provider configuration from environment variables.
func loadProviderConfig() ProviderConfig {
	return ProviderConfig{
		InfisicalClientID:     os.Getenv("INFISICAL_CLIENT_ID"),
		InfisicalClientSecret: os.Getenv("INFISICAL_CLIENT_SECRET"),
		InfisicalProjectID:    os.Getenv("INFISICAL_PROJECT_ID"),
		InfisicalSiteURL:      os.Getenv("INFISICAL_SITE_URL"),
		UseEnvFallback:        os.Getenv("USE_ENV_SECRETS") == "true",
	}
}

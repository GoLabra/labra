# Secrets Management with Infisical

Labra uses [Infisical](https://infisical.com) for secure secrets management. This document covers setup, configuration, and best practices.

## Overview

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      Infisical Cloud                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │     Dev     │  │   Staging   │  │    Prod     │              │
│  │ Environment │  │ Environment │  │ Environment │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
│         │                │                │                      │
│         └────────────────┴────────────────┘                      │
│                          │                                       │
│              Machine Identity Auth                               │
│              (IP Whitelisted)                                    │
└──────────────────────────┬──────────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
    ┌─────────▼─────────┐     ┌─────────▼─────────┐
    │  Local Dev        │     │  Production       │
    │  (CLI or SDK)     │     │  (SDK only)       │
    │                   │     │  IP Whitelisted   │
    └───────────────────┘     └───────────────────┘
```

### What's Stored Where?

| Type | Storage | Examples |
|------|---------|----------|
| **Secrets** | Infisical | `DSN`, `SECRET_KEY`, `CENTRIFUGO_API_KEY` |
| **Configuration** | Environment Variables | `DB_DIALECT`, `SERVER_PORT`, `FILE_STORAGE_PATH` |

## Initial Setup

### 1. Create Infisical Account

1. Go to [app.infisical.com](https://app.infisical.com) and create an account
2. Create a new project for Labra

### 2. Configure Environments

Create three environments in your Infisical project:
- `dev` - Local development
- `staging` - Staging/testing environment
- `prod` - Production environment

### 3. Add Secrets to Each Environment

In each environment, create the following secrets:

| Secret Key | Description |
|------------|-------------|
| `DSN` | Database connection string (e.g., `postgres://user:pass@host:5432/dbname?sslmode=disable`) |
| `SECRET_KEY` | JWT signing key (minimum 16 characters, recommended 32+) |
| `CENTRIFUGO_API_KEY` | API key for Centrifugo real-time server |

### 4. Create Machine Identity

Machine Identities are the recommended authentication method for server workloads.

1. Go to **Project Settings** → **Machine Identities**
2. Click **Create Machine Identity**
3. Name it (e.g., `labra-prod-server`)
4. Select **Universal Auth** as the authentication method
5. **Critical:** Add IP restrictions for production:
   - Add your Digital Ocean droplet's public IP(s)
   - This ensures credentials are useless if leaked
6. Copy the **Client ID** and **Client Secret**

### 5. Configure Access

Grant the Machine Identity access to your project environments:

1. Go to **Project Settings** → **Access Control**
2. Add the Machine Identity
3. Grant **Read** access to the appropriate environment(s)

## Local Development

### Option A: Infisical CLI (Recommended)

The CLI injects secrets as environment variables without storing them on disk.

1. Install the CLI:
   ```bash
   # macOS
   brew install infisical/get-cli/infisical

   # Other systems: https://infisical.com/docs/cli/overview
   ```

2. Login:
   ```bash
   infisical login
   ```

3. Run your app with secrets injected:
   ```bash
   cd resources/app
   infisical run --env=dev -- go run main.go
   ```

### Option B: SDK Authentication

Configure the SDK to fetch secrets at startup:

1. Copy `.env.example` to `.env`
2. Fill in your Infisical credentials:
   ```env
   INFISICAL_CLIENT_ID=your-client-id
   INFISICAL_CLIENT_SECRET=your-client-secret
   INFISICAL_PROJECT_ID=your-project-id
   APP_ENVIRONMENT=dev
   ```

3. Run normally:
   ```bash
   cd resources/app
   go run main.go
   ```

### Option C: Environment Fallback (Testing Only)

For quick testing without Infisical:

1. Leave `INFISICAL_CLIENT_ID` empty in your `.env`
2. Set secrets directly (not recommended for production):
   ```env
   DSN=postgres://user:pass@localhost:5432/labra
   SECRET_KEY=your-local-dev-secret-key
   CENTRIFUGO_API_KEY=your-centrifugo-key
   ```

## Production Deployment (Digital Ocean)

### Docker Configuration

Pass Infisical credentials as environment variables to your container:

```bash
docker run -d \
  -e DB_DIALECT=postgres \
  -e SERVER_PORT=4000 \
  -e CENTRIFUGO_API_ADDRESS=http://centrifugo:8000/api \
  -e APP_ENVIRONMENT=prod \
  -e INFISICAL_CLIENT_ID=your-prod-client-id \
  -e INFISICAL_CLIENT_SECRET=your-prod-client-secret \
  -e INFISICAL_PROJECT_ID=your-project-id \
  labra:latest
```

### Docker Compose Example

```yaml
version: '3.8'
services:
  labra:
    image: labra:latest
    environment:
      - DB_DIALECT=postgres
      - SERVER_PORT=4000
      - CENTRIFUGO_API_ADDRESS=http://centrifugo:8000/api
      - APP_ENVIRONMENT=prod
      - INFISICAL_CLIENT_ID=${INFISICAL_CLIENT_ID}
      - INFISICAL_CLIENT_SECRET=${INFISICAL_CLIENT_SECRET}
      - INFISICAL_PROJECT_ID=${INFISICAL_PROJECT_ID}
    ports:
      - "4000:4000"
```

### IP Whitelisting (Critical for Production)

1. Get your Digital Ocean droplet's public IP
2. In Infisical, edit your production Machine Identity
3. Add the IP to the allowlist
4. Save changes

This ensures that even if your Infisical credentials are exposed, they cannot be used from any other IP.

## Security Best Practices

### 1. Environment Isolation

- Create separate Machine Identities for each environment
- Never reuse production credentials in development
- Use IP whitelisting for production identities

### 2. Secret Rotation

To rotate secrets without downtime:

1. Update the secret value in Infisical
2. Restart your application (it will fetch the new value)
3. No code changes or redeployment needed

### 3. Audit Logging

Infisical provides audit logs for all secret access:

1. Go to **Project Settings** → **Audit Logs**
2. Monitor for unusual access patterns
3. Set up alerts for production access

### 4. Access Control

- Use the principle of least privilege
- Grant only **Read** access to application identities
- Reserve **Write** access for administrators

## Troubleshooting

### "Failed to authenticate with Infisical"

1. Verify `INFISICAL_CLIENT_ID` and `INFISICAL_CLIENT_SECRET` are correct
2. Check if the Machine Identity is active (not revoked)
3. Verify IP whitelisting allows your current IP
4. Check Infisical's status page for outages

### "Secret not found in environment"

1. Verify the secret exists in the correct environment
2. Check the secret key matches exactly (case-sensitive)
3. Ensure the Machine Identity has access to that environment

### "Using env fallback" (when you expected Infisical)

This happens when `INFISICAL_CLIENT_ID` is empty. Set it to use Infisical.

### Local Development Issues

If the CLI isn't working:
```bash
# Check login status
infisical user

# Re-authenticate
infisical login

# Verify project access
infisical secrets --env=dev
```

## Migration from .env Secrets

If you're migrating from storing secrets in `.env`:

1. Add your secrets to Infisical (all three environments)
2. Create Machine Identities for each environment
3. Update your `.env` to use Infisical credentials
4. Remove the raw secrets from `.env`
5. Test locally with `APP_ENVIRONMENT=dev`
6. Deploy with updated configuration

## References

- [Infisical Documentation](https://infisical.com/docs)
- [Go SDK Reference](https://infisical.com/docs/sdks/languages/go)
- [Machine Identity Authentication](https://infisical.com/docs/documentation/platform/identities/machine-identities)
- [CLI Installation](https://infisical.com/docs/cli/overview)


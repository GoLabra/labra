# Use the official Golang image as the base
FROM golang:1.23.0-alpine3.20

# Set the working directory inside the container
WORKDIR /app

# ============================================
# NON-SENSITIVE CONFIGURATION (passed as ARGs)
# ============================================
ARG DB_DIALECT
ARG SERVER_PORT
ARG ENT_SCHEMA_PATH
ARG CENTRIFUGO_API_ADDRESS
ARG FILE_STORAGE_PROVIDER
ARG FILE_STORAGE_PATH
ARG SUPER_ADMIN_EMAIL
ARG APP_ENVIRONMENT

# ============================================
# INFISICAL AUTHENTICATION (for secrets)
# ============================================
# These are meta-credentials for accessing the vault - NOT the actual secrets.
# Even if exposed, they:
#   1. Only work from whitelisted IPs (configure in Infisical dashboard)
#   2. Are fully audited with every access logged
#   3. Can be instantly revoked without redeployment
ARG INFISICAL_CLIENT_ID
ARG INFISICAL_CLIENT_SECRET
ARG INFISICAL_PROJECT_ID
ARG INFISICAL_SITE_URL

# Copy the entire project
COPY . .

# Change to the app directory and run Go commands
WORKDIR /app/resources/app

# Run the required commands in one shell execution
RUN go generate

RUN go mod tidy

WORKDIR /app

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 4000

# Start the app using our script
ENTRYPOINT ["/entrypoint.sh"]

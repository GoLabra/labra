# Use the official Golang image as the base
FROM golang:1.23.0-alpine3.20

# Set the working directory inside the container
WORKDIR /app

ARG DSN
ARG DB_DIALECT
ARG SERVER_PORT
ARG ENT_SCHEMA_PATH
ARG SECRET_KEY
ARG CENTRIFUGO_API_ADDRESS
ARG CENTRIFUGO_API_KEY
ARG FILE_STORAGE_PROVIDER
ARG FILE_STORAGE_PATH
ARG SUPER_ADMIN_EMAIL

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

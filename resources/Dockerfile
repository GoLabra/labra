# Use the official Golang image as the base
FROM golang:1.23.0-alpine

# Install dependencies
RUN apk add --no-cache git bash openssh-client

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

ENV DSN=$DSN
ENV DB_DIALECT=$DB_DIALECT
ENV SERVER_PORT=$SERVER_PORT
ENV ENT_SCHEMA_PATH=$ENT_SCHEMA_PATH
ENV SECRET_KEY=$SECRET_KEY
ENV CENTRIFUGO_API_ADDRESS=$CENTRIFUGO_API_ADDRESS
ENV CENTRIFUGO_API_KEY=$CENTRIFUGO_API_KEY
ENV FILE_STORAGE_PROVIDER=$FILE_STORAGE_PROVIDER
ENV FILE_STORAGE_PATH=$FILE_STORAGE_PATH

# Copy the entire project
COPY . .

# Change to the app directory and run Go commands
WORKDIR /app/app

# Run the required commands in one shell execution
RUN sed -i "/REPLACE_LABRAGO_DEVELOPMENT_API/c replace github.com\/GoLabra\/labra => ..\/..\/." go.mod  && \
    go generate && \
    go mod tidy

WORKDIR /app

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 4000

# Start the app using our script
ENTRYPOINT ["/entrypoint.sh"]

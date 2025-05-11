# Use the official Golang image as the base
FROM golang:1.23.0-alpine

# Install dependencies
RUN apk add --no-cache git bash openssh-client

# Set the working directory inside the container
WORKDIR /app

ARG GH_PAT

# Authenticate Git with GitHub Token (avoids interactive login)
RUN git config --global url."https://${GH_PAT}@github.com/".insteadOf "https://github.com/"

# Copy the entire project
COPY . .

# Change to the app directory and run Go commands
WORKDIR /app/app

# Run the required commands in one shell execution
RUN export GOPRIVATE=github.com/GoLabra/labrago && \
    go generate && \
    go mod tidy

# Expose the server port
EXPOSE 4000

# Set the working directory to the CLI
WORKDIR /app/cli

# Command to run the application
CMD ["go", "run", "main.go", "start"]

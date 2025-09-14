# Use the official Golang image as the base
FROM golang:1.23.0-alpine

# Install dependencies
RUN apk add --no-cache git bash openssh-client

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project
COPY . .

# Copy and set up entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 4000

# Start the app using our script
ENTRYPOINT ["/entrypoint.sh"]

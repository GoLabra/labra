#!/bin/sh
set -e

echo "⏳ Waiting briefly to ensure volumes are mounted..."
sleep 1

echo "📁 Checking if schema directory exists..."
if [ ! -d "/app/app/schema" ]; then
  echo "❌ Schema directory not found at /app/app/schema"
  exit 1
fi

echo "⚙️ Running go generate..."
cd /app/app
go generate

echo "📦 Running go mod tidy..."
go mod tidy

echo "🚀 Starting Go application..."
cd /app/cli
exec go run main.go start

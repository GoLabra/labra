#!/bin/sh

echo "⏳ Waiting briefly to ensure volumes are mounted..."
sleep 1

echo "📁 Checking if schema directory exists..."
if [ ! -d "/app/resources/app/schema" ]; then
  echo "❌ Schema directory not found at /app/resources/app/schema"
  exit 1
fi

echo "🚀 Starting Go application..."
cd /app/resources/cli
exec go run main.go start

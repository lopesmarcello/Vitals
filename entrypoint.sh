#!/bin/sh

# Generate Go files from templates
echo "Generating Go files from templates..."
templ generate

# Start the server
echo "Starting server..."
go run ./cmd/server/main.go

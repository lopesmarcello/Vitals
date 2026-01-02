#!/bin/sh

# Function to clean up generated files
cleanup() {
    echo "Cleaning up generated template files..."
    find . -name "*_templ.go" -type f -delete
    echo "Cleanup complete."
}

# Trap the exit signal to run the cleanup function
trap cleanup EXIT

# Generate Go files from templates
echo "Generating Go files from templates..."
templ generate

# Start the server
echo "Starting server..."
go run ./cmd/server/main.go

#!/bin/bash

# Create bin directory if it doesn't exist
mkdir -p bin

# Check for production flag
if [[ "$1" == "-p" ]] || [[ "$1" == "--prod" ]]; then
    # Build for Linux (production)
    echo "Building for Linux (production)..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags "-s -w" -o bin/ajudge
    echo "Build complete: bin/ajudge (Linux AMD64)"
else
    # Build for local system
    echo "Building for local system..."
    go build -ldflags "-s -w" -o ./ajudge
    echo "Build complete: ./ajudge ($(go env GOOS)/$(go env GOARCH))"
fi

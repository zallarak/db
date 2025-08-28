#!/bin/bash

# Kill any existing server
lsof -ti:8081 | xargs kill -9 2>/dev/null || true

# Start API server in background
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/dbxyz?sslmode=disable"
export PORT=8081
cd api && go run cmd/server/main.go &
SERVER_PID=$!

# Wait for server to start
echo "Starting API server..."
sleep 3

# Test CLI commands
echo "🧪 Testing db.xyz CLI..."
echo

cd ../cli

# Build CLI
echo "Building CLI..."
go build -o dbx cmd/dbx/main.go
echo

# Test CLI commands
echo "1. CLI Help"
./dbx --help
echo

echo "2. Test Auth Commands"
echo "Register user:"
./dbx auth register

echo
echo "3. Test Org Commands"
echo "List orgs (should fail - not logged in):"
./dbx org list
echo

# Cleanup
kill $SERVER_PID 2>/dev/null
echo "✅ CLI testing complete"
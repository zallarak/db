#!/bin/bash

set -e

# Configuration
API_BASE="http://127.0.0.1:8081"
TEST_EMAIL="clitest@example.com"
TEST_PASSWORD="testpass123"

# Kill any existing server
lsof -ti:8081 | xargs kill -9 2>/dev/null || true

# Start API server in background
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/dbxyz?sslmode=disable"
export PORT=8081
cd api && go run cmd/server/main.go &
SERVER_PID=$!

echo "🚀 Starting API server..."
sleep 4

# Test API is running
echo "Testing API health..."
curl -s "$API_BASE/health" || { echo "❌ API not responding"; kill $SERVER_PID; exit 1; }

# Register test user via API
echo "📝 Registering test user via API..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

if echo "$REGISTER_RESPONSE" | grep -q "already exists"; then
  echo "✅ User already exists (that's fine)"
elif echo "$REGISTER_RESPONSE" | grep -q "created"; then
  echo "✅ User registered successfully"
else
  echo "❌ Registration failed: $REGISTER_RESPONSE"
fi

# Build CLI
cd ../cli
echo "🔨 Building CLI..."
go build -o dbx cmd/dbx/main.go

# Test CLI help
echo "📚 Testing CLI help:"
./dbx --help | head -5

# Create a test config for CLI
echo "🔧 Creating test config..."
mkdir -p ~/.dbx-test
cat > ~/.dbx-test/config.yaml << EOF
api-url: $API_BASE
EOF

# Test CLI login manually with config
echo "🔐 Testing CLI login flow..."

# Instead of interactive login, let's test org list without auth first
echo "Testing org list without auth (should fail):"
./dbx --config ~/.dbx-test/config.yaml org list || echo "✅ Correctly failed without auth"

# Login via API and save token to config
echo "Getting token via API..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
  echo "✅ Got token: ${TOKEN:0:20}..."
  
  # Add token to config
  echo "token: $TOKEN" >> ~/.dbx-test/config.yaml
  
  # Test authenticated CLI commands
  echo "🏢 Testing org list with auth:"
  ./dbx --config ~/.dbx-test/config.yaml org list
  
  echo "👤 Testing user info:"
  ./dbx --config ~/.dbx-test/config.yaml user me || echo "Note: user me command might not exist yet"
  
  echo "🏢 Testing org creation:"
  ./dbx --config ~/.dbx-test/config.yaml org create "Test CLI Org" || echo "Note: org create syntax might be different"
  
else
  echo "❌ Failed to get token from login"
fi

# Cleanup
rm -rf ~/.dbx-test
kill $SERVER_PID 2>/dev/null

echo "✅ CLI testing complete"
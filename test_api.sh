#!/bin/bash

BASE_URL="http://127.0.0.1:8081"

echo "🧪 Testing db.xyz API endpoints..."
echo

# Function to format JSON output
format_json() {
  if command -v jq >/dev/null 2>&1; then
    echo "$1" | jq .
  else
    echo "$1"
  fi
}

# Health check
echo "1. Health Check"
HEALTH_RESPONSE=$(curl -s "$BASE_URL/health")
format_json "$HEALTH_RESPONSE"
echo

# Register a test user
echo "2. Register User"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
format_json "$REGISTER_RESPONSE"
echo

# Login with test user
echo "3. Login User"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com", 
    "password": "password123"
  }')
format_json "$LOGIN_RESPONSE"

# Extract token for authenticated requests (works with or without jq)
if command -v jq >/dev/null 2>&1; then
  TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token // empty')
else
  TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
fi
echo "Token: ${TOKEN:0:20}..." # Show first 20 chars
echo

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
  # Get current user
  echo "4. Get Current User"
  USER_RESPONSE=$(curl -s "$BASE_URL/v1/users/me" \
    -H "Authorization: Bearer $TOKEN")
  format_json "$USER_RESPONSE"
  echo

  # List organizations (should be empty initially)
  echo "5. List Organizations"
  ORGS_RESPONSE=$(curl -s "$BASE_URL/v1/orgs" \
    -H "Authorization: Bearer $TOKEN")
  format_json "$ORGS_RESPONSE"
  echo

  # Create an organization
  echo "6. Create Organization"
  ORG_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/orgs" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Test Org"
    }')
  format_json "$ORG_RESPONSE"
  echo

  # List organizations again
  echo "7. List Organizations (after creation)"
  ORGS_AFTER_RESPONSE=$(curl -s "$BASE_URL/v1/orgs" \
    -H "Authorization: Bearer $TOKEN")
  format_json "$ORGS_AFTER_RESPONSE"
  echo
else
  echo "❌ Login failed, skipping authenticated tests"
fi

echo "✅ API testing complete"
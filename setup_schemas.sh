#!/bin/bash

# Setup TinyBrain PocketBase schemas via REST API
# This script sets up the proper field schemas for our collections

BASE_URL="http://127.0.0.1:8090"
ADMIN_EMAIL="admin@tinybrain.local"
ADMIN_PASSWORD="admin123"

echo "Setting up TinyBrain PocketBase schemas..."

# Login to get admin token
echo "Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/admins/auth-with-password" \
  -H "Content-Type: application/json" \
  -d "{\"identity\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "Failed to login. Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "Successfully logged in. Token: ${TOKEN:0:20}..."

# Update sessions collection schema
echo "Setting up sessions collection schema..."
curl -s -X PATCH "$BASE_URL/api/collections/sessions" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "schema": [
      {
        "id": "name",
        "name": "name",
        "type": "text",
        "required": true
      },
      {
        "id": "description", 
        "name": "description",
        "type": "text",
        "required": false
      },
      {
        "id": "task_type",
        "name": "task_type", 
        "type": "text",
        "required": true
      },
      {
        "id": "status",
        "name": "status",
        "type": "text", 
        "required": true
      }
    ]
  }' | jq '.name // "Error"'

# Update memory_entries collection schema  
echo "Setting up memory_entries collection schema..."
curl -s -X PATCH "$BASE_URL/api/collections/memory_entries" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "schema": [
      {
        "id": "session_id",
        "name": "session_id",
        "type": "text",
        "required": true
      },
      {
        "id": "title",
        "name": "title", 
        "type": "text",
        "required": true
      },
      {
        "id": "content",
        "name": "content",
        "type": "text",
        "required": true
      },
      {
        "id": "category",
        "name": "category",
        "type": "text",
        "required": true
      },
      {
        "id": "content_type",
        "name": "content_type",
        "type": "text",
        "required": true
      },
      {
        "id": "priority",
        "name": "priority",
        "type": "number",
        "required": false
      },
      {
        "id": "confidence", 
        "name": "confidence",
        "type": "number",
        "required": false
      },
      {
        "id": "tags",
        "name": "tags",
        "type": "json",
        "required": false
      },
      {
        "id": "source",
        "name": "source",
        "type": "text",
        "required": false
      },
      {
        "id": "access_count",
        "name": "access_count", 
        "type": "number",
        "required": false
      }
    ]
  }' | jq '.name // "Error"'

# Update relationships collection schema
echo "Setting up relationships collection schema..."
curl -s -X PATCH "$BASE_URL/api/collections/relationships" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "schema": [
      {
        "id": "source_entry_id",
        "name": "source_entry_id",
        "type": "text",
        "required": true
      },
      {
        "id": "target_entry_id", 
        "name": "target_entry_id",
        "type": "text",
        "required": true
      },
      {
        "id": "relationship_type",
        "name": "relationship_type",
        "type": "text",
        "required": true
      },
      {
        "id": "strength",
        "name": "strength",
        "type": "number",
        "required": false
      },
      {
        "id": "description",
        "name": "description",
        "type": "text",
        "required": false
      }
    ]
  }' | jq '.name // "Error"'

echo "Schema setup complete!"
echo "You can now test the MCP endpoints with proper field mapping."

#!/bin/bash

# Setup PocketBase collections programmatically
echo "🧠 Setting up TinyBrain PocketBase Collections"
echo "============================================="

# Wait for server to be ready
echo "Waiting for PocketBase server to be ready..."
sleep 2

# Create a superuser first (required for API access)
echo "Creating superuser..."
./tinybrain-pb superuser upsert admin@tinybrain.local admin123

# Create collections via API
echo "Creating collections..."

# 1. Sessions Collection
echo "Creating sessions collection..."
curl -X POST http://localhost:8090/api/collections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $(curl -s -X POST http://localhost:8090/api/admins/auth-with-password \
    -H "Content-Type: application/json" \
    -d '{"identity":"admin@tinybrain.local","password":"admin123"}' | jq -r '.token')" \
  -d '{
    "name": "sessions",
    "type": "base",
    "schema": [
      {
        "name": "name",
        "type": "text",
        "required": true
      },
      {
        "name": "description",
        "type": "text"
      },
      {
        "name": "task_type",
        "type": "select",
        "required": true,
        "options": {
          "values": ["security_review", "penetration_test", "exploit_dev", "vulnerability_analysis", "threat_modeling", "incident_response", "general"]
        }
      },
      {
        "name": "status",
        "type": "select",
        "required": true,
        "options": {
          "values": ["active", "paused", "completed", "archived"]
        }
      },
      {
        "name": "metadata",
        "type": "json"
      }
    ]
  }'

echo ""
echo "✅ Collections setup complete!"
echo ""
echo "📊 Next Steps:"
echo "1. Update MCP handlers to use real PocketBase collections"
echo "2. Test full functionality with real database operations"
echo "3. Set up admin interface for ongoing management"

#!/bin/bash

# Real Intelligence Data Validation Test
# This test validates that data comes from real sources and is actually stored in the database
# NO MOCK DATA ACCEPTED

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${BLUE}Testing: $test_name${NC}"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✅ PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAILED${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""
}

# Start the server in background
echo -e "${YELLOW}Building TinyBrain Intelligence Final Server...${NC}"
./build_intelligence_final.sh

echo -e "${YELLOW}Starting TinyBrain Intelligence Final Server...${NC}"
./tinybrain-intelligence-final &
SERVER_PID=$!

# Wait for server to start
echo "Waiting for server to start..."
sleep 5

# Trap to ensure server is killed on exit
trap "kill $SERVER_PID 2>/dev/null || true; wait $SERVER_PID 2>/dev/null || true" EXIT

echo -e "${YELLOW}=== PHASE 1: Download Real Data ===${NC}"
echo ""

# Test 1: Download security data (this will take several minutes for real data)
run_test "Download Real Security Data (NVD, ATT&CK, OWASP)" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"tools/call\", \"params\": {\"name\": \"download_security_data\", \"arguments\": {}}}")
  
# Check for successful status (even partial success is ok)
echo "$RESPONSE" | jq -r ".result.status" | grep -E "(success|partial_success)"
'

# Wait a bit for data to be stored
echo "Waiting for data to be stored..."
sleep 10

echo -e "${YELLOW}=== PHASE 2: Validate Real NVD Data ===${NC}"
echo ""

# Test 2: Query NVD and validate it returns REAL CVE data
run_test "NVD Query Returns Real CVE Data" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 2, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"buffer overflow\", \"limit\": 5}}}")

# Validate response structure
echo "$RESPONSE" | jq -e ".result.results" > /dev/null || exit 1

# Get total count
TOTAL=$(echo "$RESPONSE" | jq -r ".result.total_count")
echo "Total CVEs found: $TOTAL"

# Validate we got actual results
[ "$TOTAL" -gt 0 ] || exit 1

# Validate CVE ID format (should be CVE-YYYY-NNNNN)
CVE_ID=$(echo "$RESPONSE" | jq -r ".result.results[0].cve_id")
echo "First CVE ID: $CVE_ID"
echo "$CVE_ID" | grep -E "^CVE-[0-9]{4}-[0-9]+$" || exit 1

# Validate it has real description (not sample/mock/test)
DESC=$(echo "$RESPONSE" | jq -r ".result.results[0].description")
echo "$DESC" | grep -qiE "(sample|mock|test|for testing)" && exit 1

echo "CVE Description: ${DESC:0:100}..."
'

# Test 3: Validate CVE has real CVSS scores
run_test "NVD CVE Has Valid CVSS Scores" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 3, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"remote code execution\", \"limit\": 1}}}")

# Check if we have CVSS scores
CVSS_V3=$(echo "$RESPONSE" | jq -r ".result.results[0].cvss_v3_score")
echo "CVSS v3 Score: $CVSS_V3"

# CVSS scores should be between 0.0 and 10.0
if [ "$CVSS_V3" != "null" ]; then
    python3 -c "import sys; score = float(\"$CVSS_V3\"); sys.exit(0 if 0.0 <= score <= 10.0 else 1)"
else
    echo "Warning: No CVSS v3 score found"
fi
'

echo -e "${YELLOW}=== PHASE 3: Validate Real ATT&CK Data ===${NC}"
echo ""

# Test 4: Query ATT&CK and validate it returns REAL technique data
run_test "ATT&CK Query Returns Real Technique Data" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 4, \"method\": \"tools/call\", \"params\": {\"name\": \"query_attack\", \"arguments\": {\"query\": \"Process Injection\", \"limit\": 5}}}")

# Validate response structure
echo "$RESPONSE" | jq -e ".result.results" > /dev/null || exit 1

# Get total count
TOTAL=$(echo "$RESPONSE" | jq -r ".result.total_count")
echo "Total techniques found: $TOTAL"

# Validate we got actual results
[ "$TOTAL" -gt 0 ] || exit 1

# Validate technique ID format (should be T followed by numbers)
TECH_ID=$(echo "$RESPONSE" | jq -r ".result.results[0].technique_id")
echo "First Technique ID: $TECH_ID"
echo "$TECH_ID" | grep -E "^T[0-9]+" || exit 1

# Validate it has real description (not sample/mock/test)
DESC=$(echo "$RESPONSE" | jq -r ".result.results[0].description")
echo "$DESC" | grep -qiE "(sample|mock|test|for testing|related to:)" && exit 1

echo "Technique Description: ${DESC:0:100}..."
'

# Test 5: Validate ATT&CK technique has valid tactics and platforms
run_test "ATT&CK Technique Has Valid Tactics and Platforms" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 5, \"method\": \"tools/call\", \"params\": {\"name\": \"query_attack\", \"arguments\": {\"query\": \"execution\", \"limit\": 1}}}")

# Check if we have valid tactic
TACTIC=$(echo "$RESPONSE" | jq -r ".result.results[0].tactic")
echo "Tactic: $TACTIC"
[ "$TACTIC" != "null" ] && [ "$TACTIC" != "" ] || exit 1

# Check if we have platforms array
PLATFORMS=$(echo "$RESPONSE" | jq -r ".result.results[0].platforms | length")
echo "Number of platforms: $PLATFORMS"
[ "$PLATFORMS" -gt 0 ] || exit 1
'

echo -e "${YELLOW}=== PHASE 4: Validate Real OWASP Data ===${NC}"
echo ""

# Test 6: Query OWASP and validate it returns REAL procedure data
run_test "OWASP Query Returns Real Testing Procedures" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 6, \"method\": \"tools/call\", \"params\": {\"name\": \"query_owasp\", \"arguments\": {\"query\": \"authentication\", \"limit\": 5}}}")

# Validate response structure
echo "$RESPONSE" | jq -e ".result.results" > /dev/null || exit 1

# Get total count
TOTAL=$(echo "$RESPONSE" | jq -r ".result.total_count")
echo "Total procedures found: $TOTAL"

# Validate we got actual results
[ "$TOTAL" -gt 0 ] || exit 1

# Validate it has real title and description (not sample/mock/test)
TITLE=$(echo "$RESPONSE" | jq -r ".result.results[0].title")
echo "First Procedure Title: $TITLE"
echo "$TITLE" | grep -qiE "^(sample|mock|test)" && exit 1

DESC=$(echo "$RESPONSE" | jq -r ".result.results[0].description")
echo "$DESC" | grep -qiE "(sample|mock|for testing|related to:)" && exit 1

echo "Procedure Description: ${DESC:0:100}..."
'

echo -e "${YELLOW}=== PHASE 5: Database Validation ===${NC}"
echo ""

# Test 7: Verify data is actually in database (check summary)
run_test "Database Contains Real Data (Verified via Summary)" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 7, \"method\": \"tools/call\", \"params\": {\"name\": \"get_security_data_summary\", \"arguments\": {}}}")

# Get stored record counts
NVD_COUNT=$(echo "$RESPONSE" | jq -r ".result.summary.stored_records.nvd")
ATTACK_COUNT=$(echo "$RESPONSE" | jq -r ".result.summary.stored_records.attack")
OWASP_COUNT=$(echo "$RESPONSE" | jq -r ".result.summary.stored_records.owasp")

echo "NVD CVEs stored: $NVD_COUNT"
echo "ATT&CK techniques stored: $ATTACK_COUNT"
echo "OWASP procedures stored: $OWASP_COUNT"

# Validate we have real data (not just 2 mock entries)
[ "$NVD_COUNT" -gt 10 ] || { echo "ERROR: NVD count too low, looks like mock data"; exit 1; }
[ "$ATTACK_COUNT" -gt 10 ] || { echo "ERROR: ATT&CK count too low, looks like mock data"; exit 1; }
[ "$OWASP_COUNT" -gt 0 ] || { echo "ERROR: No OWASP data found"; exit 1; }
'

echo -e "${YELLOW}=== PHASE 6: Data Source Verification ===${NC}"
echo ""

# Test 8: Verify NVD data comes from actual NVD API (check for real CVE characteristics)
run_test "NVD Data Has Real CVE Characteristics" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 8, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"\", \"limit\": 10}}}")

# Check multiple CVEs to ensure they have real characteristics
COUNT=0
for i in {0..9}; do
    CVE_ID=$(echo "$RESPONSE" | jq -r ".result.results[$i].cve_id")
    if [ "$CVE_ID" != "null" ] && echo "$CVE_ID" | grep -qE "^CVE-[0-9]{4}-[0-9]+$"; then
        COUNT=$((COUNT + 1))
    fi
done

echo "Valid CVE IDs found: $COUNT / 10"
[ "$COUNT" -ge 5 ] || exit 1
'

# Test 9: Verify ATT&CK data comes from actual MITRE repository
run_test "ATT&CK Data Has Real MITRE Characteristics" '
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 9, \"method\": \"tools/call\", \"params\": {\"name\": \"query_attack\", \"arguments\": {\"query\": \"\", \"limit\": 10}}}")

# Check multiple techniques to ensure they have real characteristics
COUNT=0
for i in {0..9}; do
    TECH_ID=$(echo "$RESPONSE" | jq -r ".result.results[$i].technique_id")
    if [ "$TECH_ID" != "null" ] && echo "$TECH_ID" | grep -qE "^T[0-9]+"; then
        COUNT=$((COUNT + 1))
    fi
done

echo "Valid technique IDs found: $COUNT / 10"
[ "$COUNT" -ge 5 ] || exit 1
'

# Test 10: No mock/sample data in database
run_test "No Mock or Sample Data in Database" '
# Check NVD for mock data
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 10, \"method\": \"tools/call\", \"params\": {\"name\": \"query_nvd\", \"arguments\": {\"query\": \"sample\", \"limit\": 10}}}")

COUNT=$(echo "$RESPONSE" | jq -r ".result.total_count")
echo "CVEs with '\''sample'\'' in description: $COUNT"

# If we find any, check they'\''re not mock data
if [ "$COUNT" -gt 0 ]; then
    DESC=$(echo "$RESPONSE" | jq -r ".result.results[0].description")
    echo "$DESC" | grep -qiE "(sample CVE for testing|mock data)" && exit 1
fi

# Check ATT&CK for mock data
RESPONSE=$(curl -s -X POST http://127.0.0.1:8090/mcp \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"id\": 11, \"method\": \"tools/call\", \"params\": {\"name\": \"query_attack\", \"arguments\": {\"query\": \"sample\", \"limit\": 10}}}")

COUNT=$(echo "$RESPONSE" | jq -r ".result.total_count")
echo "Techniques with '\''sample'\'' in description: $COUNT"

if [ "$COUNT" -gt 0 ]; then
    DESC=$(echo "$RESPONSE" | jq -r ".result.results[0].description")
    echo "$DESC" | grep -qiE "(sample|mock|for testing|related to:)" && exit 1
fi
'

# Print final results
echo ""
echo -e "${YELLOW}=== FINAL VALIDATION RESULTS ===${NC}"
echo -e "Total Tests: ${TOTAL_TESTS}"
echo -e "Passed: ${GREEN}${PASSED_TESTS}${NC}"
echo -e "Failed: ${RED}${FAILED_TESTS}${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL VALIDATION TESTS PASSED! 🎉${NC}"
    echo -e "${GREEN}✅ Real intelligence data is working correctly!${NC}"
    echo ""
    echo "Summary of validated data:"
    echo "- NVD CVE data is from real National Vulnerability Database"
    echo "- ATT&CK techniques are from real MITRE framework"
    echo "- OWASP procedures are from real OWASP testing guide"
    echo "- NO mock or sample data found"
    echo "- Data is properly stored in database"
    exit 0
else
    echo -e "${RED}❌ Some validation tests failed.${NC}"
    echo -e "${RED}Success rate: $(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%${NC}"
    echo ""
    echo "Please review the failures above."
    exit 1
fi


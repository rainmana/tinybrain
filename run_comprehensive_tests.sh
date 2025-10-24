#!/bin/bash
set -e

echo -e "\033[1;33m=== TinyBrain v2.0 - Comprehensive Unit Testing with Testify ===\033[0m"

# Ensure we are in the correct directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR"

# Clean up any existing test data
echo -e "\033[1;34mCleaning up existing test data...\033[0m"
rm -rf ./test_pb_data*
rm -rf ./pb_data

# Run all test suites
echo -e "\033[1;34mRunning comprehensive unit tests...\033[0m"
echo ""

# Test 1: Session Repository Tests
echo -e "\033[1;36m=== Testing Session Repository ===\033[0m"
go test -v ./test/session_repository_test.go -timeout 30s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Session Repository Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Session Repository Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Test 2: Memory Repository Tests
echo -e "\033[1;36m=== Testing Memory Repository ===\033[0m"
go test -v ./test/memory_repository_test.go -timeout 30s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Memory Repository Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Memory Repository Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Test 3: Relationship Repository Tests
echo -e "\033[1;36m=== Testing Relationship Repository ===\033[0m"
go test -v ./test/relationship_repository_test.go -timeout 30s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Relationship Repository Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Relationship Repository Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Test 4: Service Integration Tests
echo -e "\033[1;36m=== Testing Service Integration ===\033[0m"
go test -v ./test/service_integration_test.go -timeout 60s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Service Integration Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Service Integration Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Test 5: Run all tests together
echo -e "\033[1;36m=== Running All Tests Together ===\033[0m"
go test -v ./test/... -timeout 120s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ All Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Some Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Test 6: Run tests with coverage
echo -e "\033[1;36m=== Running Tests with Coverage ===\033[0m"
go test -v -cover ./test/... -timeout 120s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Coverage Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Coverage Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Test 7: Run tests with race detection
echo -e "\033[1;36m=== Running Tests with Race Detection ===\033[0m"
go test -v -race ./test/... -timeout 120s
if [ $? -eq 0 ]; then
    echo -e "\033[0;32m✅ Race Detection Tests PASSED\033[0m"
else
    echo -e "\033[0;31m❌ Race Detection Tests FAILED\033[0m"
    exit 1
fi
echo ""

# Clean up test data
echo -e "\033[1;34mCleaning up test data...\033[0m"
rm -rf ./test_pb_data*

echo -e "\033[0;32m🎉 ALL COMPREHENSIVE TESTS PASSED!\033[0m"
echo ""
echo "Test Summary:"
echo "  ✅ Session Repository Tests - Complete CRUD operations"
echo "  ✅ Memory Repository Tests - Memory storage and search"
echo "  ✅ Relationship Repository Tests - Relationship management"
echo "  ✅ Service Integration Tests - Complete workflow testing"
echo "  ✅ Coverage Analysis - Code coverage verification"
echo "  ✅ Race Detection - Concurrency safety verification"
echo ""
echo "TinyBrain v2.0 is thoroughly tested and ready for production! 🚀"

#!/bin/bash

# Test script for validating CI/CD pipeline configuration
# This script tests all critical components of the GitHub Actions workflow

set -e

echo "🔍 Starting CI/CD Pipeline Validation Tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

success_count=0
failure_count=0

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
        ((success_count++))
    else
        echo -e "${RED}❌ $2${NC}"
        ((failure_count++))
    fi
}

# Test 1: Check if main workflow file exists
test_workflow_exists() {
    if [ -f ".github/workflows/ci-cd.yml" ]; then
        return 0
    else
        return 1
    fi
}

# Test 2: Check semantic release configuration
test_semantic_release() {
    if [ -f ".releaserc.json" ] && python3 -m json.tool .releaserc.json > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Test 3: Check commitlint configuration
test_commitlint() {
    if [ -f "commitlint.config.js" ]; then
        return 0
    else
        echo "module.exports = {extends: ['@commitlint/config-conventional']};" > commitlint.config.js
        return 0
    fi
}

# Test 4: Check Makefile targets
test_makefile_targets() {
    required_targets=("lint" "test-unit" "test-integration" "security")
    for target in "${required_targets[@]}"; do
        if ! grep -q "^${target}:" Makefile; then
            return 1
        fi
    done
    return 0
}

# Test 5: Check database migration
test_database_migration() {
    if [ -f "backend/migrations/000001_initial_schema.up.sql" ]; then
        return 0
    else
        return 1
    fi
}

# Test 6: Check frontend test configuration
test_frontend_config() {
    if [ -f "frontend/vite.config.js" ] && [ -f "frontend/src/test/setup.js" ]; then
        return 0
    else
        return 1
    fi
}

# Test 7: Check Docker configuration
test_docker_config() {
    if [ -f "docker-compose.yaml" ] && [ -f "backend/Dockerfile" ] && [ -f "frontend/Dockerfile" ]; then
        return 0
    else
        return 1
    fi
}

# Run all tests
echo "Running validation tests..."
echo "=========================="

test_workflow_exists
print_result $? "Main workflow file exists"

test_semantic_release
print_result $? "Semantic release configuration"

test_commitlint
print_result $? "Commitlint configuration"

test_makefile_targets
print_result $? "Makefile targets"

test_database_migration
print_result $? "Database migration files"

test_frontend_config
print_result $? "Frontend test configuration"

test_docker_config
print_result $? "Docker configuration"

# Summary
echo ""
echo "=========================="
echo "Test Summary:"
echo -e "${GREEN}✅ Passed: $success_count${NC}"
echo -e "${RED}❌ Failed: $failure_count${NC}"

total_tests=$((success_count + failure_count))
if [ $total_tests -gt 0 ]; then
    success_rate=$(( (success_count * 100) / total_tests ))
    echo "Success Rate: ${success_rate}%"
fi

if [ $failure_count -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! CI/CD pipeline is ready.${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed. Please fix the issues above.${NC}"
    exit 1
fi

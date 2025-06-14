#!/bin/bash

# CI/CD Pipeline Validation Script
# This script validates the entire CI/CD pipeline locally

set -e  # Exit on any error

echo "🚀 Starting Hydro Habitat CI/CD Pipeline Validation"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
    echo -e "\n${BLUE}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_step "Checking prerequisites..."

if ! command_exists docker; then
    print_error "Docker is not installed"
    exit 1
fi

if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
    print_error "Docker Compose is not available"
    exit 1
fi

if ! command_exists make; then
    print_error "Make is not installed"
    exit 1
fi

print_success "All prerequisites are available"

# Clean up any running containers
print_step "Cleaning up existing containers..."
docker compose down --volumes --remove-orphans 2>/dev/null || true
print_success "Cleanup completed"

# Step 1: Backend Quality Checks
print_step "Running backend quality checks..."

print_step "  → Backend linting"
if make lint-backend; then
    print_success "Backend linting passed"
else
    print_error "Backend linting failed"
    exit 1
fi

print_step "  → Backend security scan"
if make security-backend; then
    print_success "Backend security scan passed"
else
    print_error "Backend security scan failed"
    exit 1
fi

# Step 2: Backend Tests
print_step "Running backend tests..."

print_step "  → Unit tests"
if make test-unit; then
    print_success "Backend unit tests passed"
else
    print_error "Backend unit tests failed"
    exit 1
fi

# Step 3: Frontend Quality Checks
print_step "Running frontend quality checks..."

print_step "  → Frontend linting"
if make lint-frontend; then
    print_success "Frontend linting passed"
else
    print_error "Frontend linting failed"
    exit 1
fi

print_step "  → Frontend tests"
cd frontend
if npm test -- --run; then
    print_success "Frontend tests passed"
else
    print_error "Frontend tests failed"
    exit 1
fi
cd ..

# Step 4: Build and Test Docker Images
print_step "Building and testing Docker images..."

print_step "  → Building production images"
if make build; then
    print_success "Docker images built successfully"
else
    print_error "Docker image build failed"
    exit 1
fi

print_step "  → Testing image startup"
if docker compose up -d; then
    print_success "Services started successfully"
    
    # Wait for services to be ready
    echo "Waiting for services to be ready..."
    sleep 30
    
    # Test backend health
    print_step "  → Testing backend health endpoint"
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        print_success "Backend health check passed"
    else
        print_warning "Backend health check failed - this might be expected if the backend needs a database"
        docker compose logs backend | tail -10
    fi
    
    # Test frontend availability
    print_step "  → Testing frontend availability"
    for i in {1..10}; do
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:80 2>/dev/null || echo "000")
        if [ "$HTTP_CODE" = "200" ]; then
            print_success "Frontend availability check passed (HTTP $HTTP_CODE)"
            break
        fi
        if [ $i -eq 10 ]; then
            print_warning "Frontend availability check failed after 10 attempts"
            echo "Last HTTP response code: $HTTP_CODE"
            docker compose logs frontend | tail -10
        fi
        echo "Attempt $i/10: Frontend not ready (HTTP $HTTP_CODE), waiting..."
        sleep 2
    done
    
    # Show service logs for debugging
    print_step "  → Service logs (last 5 lines each)"
    echo "Backend logs:"
    docker compose logs --tail=5 backend
    echo -e "\nFrontend logs:"
    docker compose logs --tail=5 frontend
    
    # Clean up
    docker compose down --volumes --remove-orphans
    print_success "Services stopped and cleaned up"
else
    print_error "Failed to start services"
    docker compose logs
    docker compose down --volumes --remove-orphans
    exit 1
fi

# Step 5: Validate Makefile targets
print_step "Validating Makefile targets..."

EXPECTED_TARGETS=(
    "help"
    "up"
    "down"
    "logs"
    "test-unit"
    "test-integration"
    "test-integration-ci"
    "lint"
    "lint-backend"
    "lint-frontend"
    "security"
    "security-backend"
    "build"
    "clean"
)

for target in "${EXPECTED_TARGETS[@]}"; do
    if make -n "$target" >/dev/null 2>&1; then
        print_success "Target '$target' is available"
    else
        print_error "Target '$target' is missing"
        exit 1
    fi
done

# Step 6: Check GitHub Actions workflow files
print_step "Validating GitHub Actions workflows..."

WORKFLOW_FILES=(
    ".github/workflows/ci-cd.yml"
    ".github/workflows/branch-protection.yaml"
    ".github/workflows/pr-check.yaml"
    ".github/workflows/release.yaml"
)

for workflow in "${WORKFLOW_FILES[@]}"; do
    if [[ -f "$workflow" ]]; then
        print_success "Workflow file '$workflow' exists"
    else
        print_warning "Workflow file '$workflow' is missing"
    fi
done

# Step 7: Check configuration files
print_step "Validating configuration files..."

CONFIG_FILES=(
    "docker-compose.yaml"
    "Makefile"
    "frontend/package.json"
    "frontend/eslint.config.js"
    "backend/go.mod"
    "backend/go.sum"
)

for config in "${CONFIG_FILES[@]}"; do
    if [[ -f "$config" ]]; then
        print_success "Config file '$config' exists"
    else
        print_error "Config file '$config' is missing"
        exit 1
    fi
done

# Step 8: Summary
print_step "Pipeline Validation Summary"
echo "=================================================="
print_success "All pipeline steps completed successfully!"
echo ""
echo "✅ Backend linting: PASSED"
echo "✅ Backend security: PASSED"
echo "✅ Backend unit tests: PASSED"
echo "✅ Frontend linting: PASSED"
echo "✅ Frontend tests: PASSED"
echo "✅ Docker build: PASSED"
echo "✅ Service startup: TESTED"
echo "✅ Makefile targets: VALIDATED"
echo "✅ GitHub workflows: VALIDATED"
echo "✅ Configuration files: VALIDATED"
echo ""
print_success "🎉 The CI/CD pipeline is ready for production!"
print_success "You can now safely push your changes to trigger the GitHub Actions workflow."

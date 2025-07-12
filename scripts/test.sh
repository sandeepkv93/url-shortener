#!/bin/bash

# Test Automation Script for URL Shortener
# Runs comprehensive tests across backend and frontend

set -euo pipefail

# Script Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
URL Shortener Test Automation Script

Usage: $0 [OPTIONS] [COMMAND]

Commands:
    all             Run all tests (default)
    unit            Run unit tests only
    integration     Run integration tests only
    e2e             Run end-to-end tests only
    backend         Run backend tests only
    frontend        Run frontend tests only
    coverage        Generate test coverage reports
    performance     Run performance tests
    security        Run security tests
    lint            Run linting checks
    clean           Clean test artifacts

Options:
    -h, --help              Show this help message
    -v, --verbose           Enable verbose output
    -q, --quiet             Suppress non-essential output
    -f, --fail-fast         Stop on first test failure
    -c, --coverage          Generate coverage reports
    -p, --parallel          Run tests in parallel where possible
    -w, --watch             Watch mode for continuous testing
    --ci                    CI mode (optimized for CI/CD)
    --docker                Run tests in Docker containers
    --no-cache              Disable test caching

Examples:
    $0                      # Run all tests
    $0 unit                 # Run unit tests only
    $0 backend --coverage   # Run backend tests with coverage
    $0 frontend --watch     # Run frontend tests in watch mode
    $0 --ci                 # Run tests in CI mode
    $0 --docker             # Run tests in Docker
EOF
}

# Configuration
COMMAND="all"
VERBOSE=false
QUIET=false
FAIL_FAST=false
COVERAGE=false
PARALLEL=false
WATCH=false
CI_MODE=false
DOCKER_MODE=false
NO_CACHE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        -f|--fail-fast)
            FAIL_FAST=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -p|--parallel)
            PARALLEL=true
            shift
            ;;
        -w|--watch)
            WATCH=true
            shift
            ;;
        --ci)
            CI_MODE=true
            shift
            ;;
        --docker)
            DOCKER_MODE=true
            shift
            ;;
        --no-cache)
            NO_CACHE=true
            shift
            ;;
        all|unit|integration|e2e|backend|frontend|coverage|performance|security|lint|clean)
            COMMAND="$1"
            shift
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Helper functions
run_command() {
    local cmd="$1"
    local desc="${2:-Running command}"
    
    if [[ "$QUIET" == false ]]; then
        log "$desc"
    fi
    
    if [[ "$VERBOSE" == true ]]; then
        log "Command: $cmd"
    fi
    
    if ! eval "$cmd"; then
        error "Command failed: $cmd"
    fi
}

check_dependencies() {
    log "Checking test dependencies..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        error "Go is not installed"
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        error "Node.js is not installed"
    fi
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        error "npm is not installed"
    fi
    
    # Check Docker if needed
    if [[ "$DOCKER_MODE" == true ]]; then
        if ! command -v docker &> /dev/null; then
            error "Docker is not installed"
        fi
        
        if ! command -v docker-compose &> /dev/null; then
            error "Docker Compose is not installed"
        fi
    fi
    
    success "Dependencies check passed"
}

# Backend test functions
run_backend_unit_tests() {
    log "Running backend unit tests..."
    
    cd "$PROJECT_ROOT/backend"
    
    local go_test_args=""
    
    if [[ "$VERBOSE" == true ]]; then
        go_test_args="$go_test_args -v"
    fi
    
    if [[ "$FAIL_FAST" == true ]]; then
        go_test_args="$go_test_args -failfast"
    fi
    
    if [[ "$COVERAGE" == true ]]; then
        go_test_args="$go_test_args -coverprofile=coverage.out -covermode=atomic"
    fi
    
    if [[ "$PARALLEL" == true ]]; then
        go_test_args="$go_test_args -parallel=4"
    fi
    
    if [[ "$CI_MODE" == true ]]; then
        go_test_args="$go_test_args -race -timeout=5m"
    fi
    
    # Run unit tests
    run_command "go test $go_test_args ./internal/..." "Running Go unit tests"
    
    # Generate coverage report if requested
    if [[ "$COVERAGE" == true ]]; then
        run_command "go tool cover -html=coverage.out -o coverage.html" "Generating coverage report"
        run_command "go tool cover -func=coverage.out" "Displaying coverage summary"
    fi
    
    cd "$PROJECT_ROOT"
    success "Backend unit tests completed"
}

run_backend_integration_tests() {
    log "Running backend integration tests..."
    
    cd "$PROJECT_ROOT/backend"
    
    # Set up test environment
    export GO_ENV=test
    export DATABASE_URL="postgres://test:test@localhost:5433/urlshortener_test?sslmode=disable"
    export REDIS_URL="redis://:test@localhost:6380/1"
    
    # Start test dependencies if not in Docker mode
    if [[ "$DOCKER_MODE" == false ]]; then
        start_test_dependencies
    fi
    
    local go_test_args="-tags=integration"
    
    if [[ "$VERBOSE" == true ]]; then
        go_test_args="$go_test_args -v"
    fi
    
    if [[ "$FAIL_FAST" == true ]]; then
        go_test_args="$go_test_args -failfast"
    fi
    
    if [[ "$COVERAGE" == true ]]; then
        go_test_args="$go_test_args -coverprofile=integration-coverage.out -covermode=atomic"
    fi
    
    # Run integration tests
    run_command "go test $go_test_args ./tests/integration/..." "Running Go integration tests"
    
    # Clean up test dependencies
    if [[ "$DOCKER_MODE" == false ]]; then
        stop_test_dependencies
    fi
    
    cd "$PROJECT_ROOT"
    success "Backend integration tests completed"
}

start_test_dependencies() {
    log "Starting test dependencies..."
    
    # Start PostgreSQL test instance
    if ! docker ps | grep -q "postgres-test"; then
        run_command "docker run -d --name postgres-test -p 5433:5432 -e POSTGRES_DB=urlshortener_test -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test postgres:15-alpine" "Starting test PostgreSQL"
        sleep 5
    fi
    
    # Start Redis test instance
    if ! docker ps | grep -q "redis-test"; then
        run_command "docker run -d --name redis-test -p 6380:6379 redis:7-alpine redis-server --requirepass test" "Starting test Redis"
        sleep 2
    fi
}

stop_test_dependencies() {
    log "Stopping test dependencies..."
    
    docker stop postgres-test redis-test 2>/dev/null || true
    docker rm postgres-test redis-test 2>/dev/null || true
}

# Frontend test functions
run_frontend_unit_tests() {
    log "Running frontend unit tests..."
    
    cd "$PROJECT_ROOT/frontend"
    
    local npm_test_args=""
    
    if [[ "$COVERAGE" == true ]]; then
        npm_test_args="--coverage"
    fi
    
    if [[ "$WATCH" == true ]]; then
        npm_test_args="$npm_test_args --watch"
    fi
    
    if [[ "$CI_MODE" == true ]]; then
        npm_test_args="$npm_test_args --watchAll=false --ci"
    fi
    
    # Run unit tests
    run_command "npm test $npm_test_args" "Running frontend unit tests"
    
    cd "$PROJECT_ROOT"
    success "Frontend unit tests completed"
}

run_frontend_e2e_tests() {
    log "Running frontend E2E tests..."
    
    cd "$PROJECT_ROOT/frontend"
    
    # Start application if not running
    local app_started=false
    if ! curl -sf http://localhost:3000 &> /dev/null; then
        log "Starting application for E2E tests..."
        npm run build
        npm run preview &
        APP_PID=$!
        app_started=true
        sleep 5
    fi
    
    # Run E2E tests
    if [[ "$CI_MODE" == true ]]; then
        run_command "npx playwright test --reporter=github" "Running Playwright E2E tests"
    else
        run_command "npx playwright test" "Running Playwright E2E tests"
    fi
    
    # Clean up
    if [[ "$app_started" == true ]]; then
        kill $APP_PID 2>/dev/null || true
    fi
    
    cd "$PROJECT_ROOT"
    success "Frontend E2E tests completed"
}

# Linting functions
run_backend_lint() {
    log "Running backend linting..."
    
    cd "$PROJECT_ROOT/backend"
    
    # Go fmt
    run_command "gofmt -l ." "Checking Go formatting"
    local fmt_issues=$(gofmt -l . | wc -l)
    if [[ $fmt_issues -gt 0 ]]; then
        warning "Found $fmt_issues formatting issues. Run 'gofmt -w .' to fix."
    fi
    
    # Go vet
    run_command "go vet ./..." "Running go vet"
    
    # golangci-lint if available
    if command -v golangci-lint &> /dev/null; then
        run_command "golangci-lint run" "Running golangci-lint"
    else
        warning "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi
    
    cd "$PROJECT_ROOT"
    success "Backend linting completed"
}

run_frontend_lint() {
    log "Running frontend linting..."
    
    cd "$PROJECT_ROOT/frontend"
    
    # ESLint
    run_command "npm run lint" "Running ESLint"
    
    # TypeScript check
    run_command "npx tsc --noEmit" "Running TypeScript check"
    
    # Prettier check
    if command -v prettier &> /dev/null; then
        run_command "prettier --check src/" "Checking code formatting"
    fi
    
    cd "$PROJECT_ROOT"
    success "Frontend linting completed"
}

# Security tests
run_security_tests() {
    log "Running security tests..."
    
    # Go security scan
    if command -v gosec &> /dev/null; then
        cd "$PROJECT_ROOT/backend"
        run_command "gosec ./..." "Running Go security scan"
        cd "$PROJECT_ROOT"
    else
        warning "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
    fi
    
    # npm audit
    cd "$PROJECT_ROOT/frontend"
    run_command "npm audit --audit-level=moderate" "Running npm audit"
    cd "$PROJECT_ROOT"
    
    # Docker security scan if available
    if command -v docker &> /dev/null && command -v trivy &> /dev/null; then
        run_command "trivy image urlshortener-backend:latest" "Scanning backend Docker image"
        run_command "trivy image urlshortener-frontend:latest" "Scanning frontend Docker image"
    fi
    
    success "Security tests completed"
}

# Performance tests
run_performance_tests() {
    log "Running performance tests..."
    
    # Backend performance tests
    cd "$PROJECT_ROOT/backend"
    if [[ -d "tests/performance" ]]; then
        run_command "go test -bench=. -benchmem ./tests/performance/..." "Running Go benchmark tests"
    else
        warning "No backend performance tests found"
    fi
    cd "$PROJECT_ROOT"
    
    # Frontend performance tests with Lighthouse CI
    if command -v lhci &> /dev/null; then
        cd "$PROJECT_ROOT/frontend"
        run_command "lhci autorun" "Running Lighthouse CI"
        cd "$PROJECT_ROOT"
    else
        warning "Lighthouse CI not found. Install with: npm install -g @lhci/cli"
    fi
    
    success "Performance tests completed"
}

# Docker-based testing
run_docker_tests() {
    log "Running tests in Docker..."
    
    cd "$PROJECT_ROOT"
    
    # Build test images
    run_command "docker-compose -f docker-compose.test.yml build" "Building test images"
    
    # Run tests
    run_command "docker-compose -f docker-compose.test.yml up --abort-on-container-exit" "Running tests in Docker"
    
    # Extract test results
    run_command "docker-compose -f docker-compose.test.yml logs > test-results.log" "Extracting test results"
    
    # Clean up
    run_command "docker-compose -f docker-compose.test.yml down -v" "Cleaning up test containers"
    
    success "Docker tests completed"
}

# Coverage report generation
generate_coverage_report() {
    log "Generating comprehensive coverage report..."
    
    local coverage_dir="$PROJECT_ROOT/coverage"
    mkdir -p "$coverage_dir"
    
    # Backend coverage
    if [[ -f "$PROJECT_ROOT/backend/coverage.out" ]]; then
        cd "$PROJECT_ROOT/backend"
        run_command "go tool cover -html=coverage.out -o $coverage_dir/backend-coverage.html" "Generating backend coverage HTML"
        run_command "go tool cover -func=coverage.out | tail -1" "Backend coverage summary"
        cd "$PROJECT_ROOT"
    fi
    
    # Frontend coverage
    if [[ -d "$PROJECT_ROOT/frontend/coverage" ]]; then
        cp -r "$PROJECT_ROOT/frontend/coverage" "$coverage_dir/frontend"
        log "Frontend coverage report copied to $coverage_dir/frontend"
    fi
    
    # Combined coverage report
    if command -v gocovmerge &> /dev/null; then
        run_command "gocovmerge backend/coverage.out backend/integration-coverage.out > $coverage_dir/combined-coverage.out" "Merging coverage reports"
    fi
    
    success "Coverage report generated in $coverage_dir"
}

# Clean up test artifacts
clean_test_artifacts() {
    log "Cleaning test artifacts..."
    
    # Backend
    rm -f "$PROJECT_ROOT/backend/coverage.out"
    rm -f "$PROJECT_ROOT/backend/coverage.html"
    rm -f "$PROJECT_ROOT/backend/integration-coverage.out"
    
    # Frontend
    rm -rf "$PROJECT_ROOT/frontend/coverage"
    rm -rf "$PROJECT_ROOT/frontend/test-results"
    rm -rf "$PROJECT_ROOT/frontend/playwright-report"
    
    # General
    rm -rf "$PROJECT_ROOT/coverage"
    rm -f "$PROJECT_ROOT/test-results.log"
    
    # Docker
    docker-compose -f docker-compose.test.yml down -v 2>/dev/null || true
    
    success "Test artifacts cleaned"
}

# Main test execution
main() {
    log "Starting test execution - Command: $COMMAND"
    
    # Check dependencies
    check_dependencies
    
    # Set CI mode defaults
    if [[ "$CI_MODE" == true ]]; then
        COVERAGE=true
        PARALLEL=true
        FAIL_FAST=true
    fi
    
    # Execute based on command
    case "$COMMAND" in
        "all")
            if [[ "$DOCKER_MODE" == true ]]; then
                run_docker_tests
            else
                run_backend_unit_tests
                run_backend_integration_tests
                run_frontend_unit_tests
                run_frontend_e2e_tests
                run_backend_lint
                run_frontend_lint
            fi
            ;;
        "unit")
            run_backend_unit_tests
            run_frontend_unit_tests
            ;;
        "integration")
            run_backend_integration_tests
            ;;
        "e2e")
            run_frontend_e2e_tests
            ;;
        "backend")
            run_backend_unit_tests
            run_backend_integration_tests
            run_backend_lint
            ;;
        "frontend")
            run_frontend_unit_tests
            run_frontend_e2e_tests
            run_frontend_lint
            ;;
        "coverage")
            COVERAGE=true
            run_backend_unit_tests
            run_backend_integration_tests
            run_frontend_unit_tests
            generate_coverage_report
            ;;
        "performance")
            run_performance_tests
            ;;
        "security")
            run_security_tests
            ;;
        "lint")
            run_backend_lint
            run_frontend_lint
            ;;
        "clean")
            clean_test_artifacts
            ;;
        *)
            error "Unknown command: $COMMAND"
            ;;
    esac
    
    # Generate coverage report if requested
    if [[ "$COVERAGE" == true ]] && [[ "$COMMAND" != "coverage" ]] && [[ "$COMMAND" != "clean" ]]; then
        generate_coverage_report
    fi
    
    success "Test execution completed successfully"
}

# Trap to clean up on exit
trap 'stop_test_dependencies 2>/dev/null || true' EXIT

# Run main function
main "$@"
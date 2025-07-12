#!/bin/bash

# Test Graceful Shutdown Script
# This script tests the graceful shutdown functionality of the URL Shortener service

set -e

echo "🔧 Testing Graceful Shutdown for URL Shortener Service"
echo "=================================================="

# Configuration
SERVICE_URL="http://localhost:8080"
HEALTH_ENDPOINT="${SERVICE_URL}/health"
SHUTDOWN_TEST_DURATION=5

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if service is running
check_service_health() {
    if curl -s "${HEALTH_ENDPOINT}" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Function to get service PID
get_service_pid() {
    pgrep -f "url-shortener" | head -1
}

# Function to test graceful shutdown with SIGTERM
test_sigterm_shutdown() {
    print_status "Testing graceful shutdown with SIGTERM signal..."
    
    local pid=$(get_service_pid)
    if [ -z "$pid" ]; then
        print_error "Service is not running. Please start the service first."
        return 1
    fi
    
    print_status "Service PID: $pid"
    
    # Check initial health
    if check_service_health; then
        print_success "Service is healthy before shutdown"
    else
        print_error "Service is not healthy before shutdown"
        return 1
    fi
    
    # Send SIGTERM signal
    print_status "Sending SIGTERM signal to PID $pid..."
    kill -TERM "$pid"
    
    # Monitor shutdown process
    local shutdown_start=$(date +%s)
    local timeout=30
    
    while kill -0 "$pid" 2>/dev/null; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - shutdown_start))
        
        if [ $elapsed -gt $timeout ]; then
            print_error "Graceful shutdown timed out after ${timeout} seconds"
            return 1
        fi
        
        print_status "Shutdown in progress... (${elapsed}s elapsed)"
        sleep 1
    done
    
    local shutdown_duration=$(($(date +%s) - shutdown_start))
    print_success "Graceful shutdown completed in ${shutdown_duration} seconds"
    
    # Verify service is no longer responding
    if check_service_health; then
        print_error "Service is still responding after shutdown"
        return 1
    else
        print_success "Service properly stopped responding"
    fi
    
    return 0
}

# Function to test graceful shutdown with SIGINT (Ctrl+C)
test_sigint_shutdown() {
    print_status "Testing graceful shutdown with SIGINT signal..."
    
    local pid=$(get_service_pid)
    if [ -z "$pid" ]; then
        print_error "Service is not running. Please start the service first."
        return 1
    fi
    
    print_status "Service PID: $pid"
    
    # Send SIGINT signal (same as Ctrl+C)
    print_status "Sending SIGINT signal to PID $pid..."
    kill -INT "$pid"
    
    # Monitor shutdown process (similar to SIGTERM test)
    local shutdown_start=$(date +%s)
    local timeout=30
    
    while kill -0 "$pid" 2>/dev/null; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - shutdown_start))
        
        if [ $elapsed -gt $timeout ]; then
            print_error "Graceful shutdown timed out after ${timeout} seconds"
            return 1
        fi
        
        print_status "Shutdown in progress... (${elapsed}s elapsed)"
        sleep 1
    done
    
    local shutdown_duration=$(($(date +%s) - shutdown_start))
    print_success "Graceful shutdown completed in ${shutdown_duration} seconds"
    
    return 0
}

# Function to test health endpoints during startup
test_health_endpoints() {
    print_status "Testing health endpoints..."
    
    local endpoints=("/health" "/healthz" "/ready" "/readiness" "/live" "/liveness")
    
    for endpoint in "${endpoints[@]}"; do
        local url="${SERVICE_URL}${endpoint}"
        print_status "Testing endpoint: $endpoint"
        
        if curl -s -f "$url" > /dev/null; then
            print_success "✓ $endpoint is responding"
        else
            print_error "✗ $endpoint is not responding"
        fi
    done
}

# Function to test debug endpoints (development mode)
test_debug_endpoints() {
    print_status "Testing debug endpoints (if in development mode)..."
    
    local debug_endpoints=("/debug/db-stats" "/debug/redis-info" "/debug/shutdown-status")
    
    for endpoint in "${debug_endpoints[@]}"; do
        local url="${SERVICE_URL}${endpoint}"
        print_status "Testing debug endpoint: $endpoint"
        
        if curl -s -f "$url" > /dev/null; then
            print_success "✓ $endpoint is responding"
        else
            print_warning "✗ $endpoint is not responding (may be disabled in production)"
        fi
    done
}

# Function to simulate load during shutdown
simulate_load_during_shutdown() {
    print_status "Simulating load during graceful shutdown..."
    
    # Start background requests
    for i in {1..10}; do
        curl -s "${SERVICE_URL}/health" > /dev/null &
    done
    
    local pid=$(get_service_pid)
    if [ -z "$pid" ]; then
        print_error "Service is not running"
        return 1
    fi
    
    # Send shutdown signal while requests are in flight
    kill -TERM "$pid"
    
    # Wait for shutdown
    wait
    
    print_success "Load simulation during shutdown completed"
}

# Main test function
run_tests() {
    print_status "Starting graceful shutdown tests..."
    echo
    
    # Check if service is running
    if ! check_service_health; then
        print_error "Service is not running on ${SERVICE_URL}"
        print_status "Please start the service first with: cd backend && go run cmd/server/main.go"
        exit 1
    fi
    
    print_success "Service is running and healthy"
    echo
    
    # Test health endpoints first
    test_health_endpoints
    echo
    
    # Test debug endpoints
    test_debug_endpoints
    echo
    
    # Test graceful shutdown
    case "${1:-sigterm}" in
        "sigterm")
            test_sigterm_shutdown
            ;;
        "sigint")
            test_sigint_shutdown
            ;;
        "load")
            simulate_load_during_shutdown
            ;;
        "all")
            test_sigterm_shutdown
            echo
            print_status "Waiting 5 seconds before next test..."
            sleep 5
            # Would need to restart service for additional tests
            print_warning "Please restart the service to run additional tests"
            ;;
        *)
            print_error "Unknown test type: $1"
            echo "Usage: $0 [sigterm|sigint|load|all]"
            exit 1
            ;;
    esac
}

# Cleanup function
cleanup() {
    print_status "Cleaning up any remaining processes..."
    pkill -f "url-shortener" 2>/dev/null || true
}

# Trap cleanup on script exit
trap cleanup EXIT

# Parse command line arguments
case "${1:-help}" in
    "help"|"-h"|"--help")
        echo "URL Shortener Graceful Shutdown Test Script"
        echo
        echo "Usage: $0 [test_type]"
        echo
        echo "Test types:"
        echo "  sigterm  - Test graceful shutdown with SIGTERM signal (default)"
        echo "  sigint   - Test graceful shutdown with SIGINT signal (Ctrl+C)"
        echo "  load     - Test shutdown while handling requests"
        echo "  all      - Run all tests (requires manual service restarts)"
        echo "  help     - Show this help message"
        echo
        echo "Prerequisites:"
        echo "  - Service must be running on ${SERVICE_URL}"
        echo "  - curl must be installed"
        echo "  - Service must be built and ready to run"
        echo
        echo "Example:"
        echo "  # Terminal 1: Start the service"
        echo "  cd backend && go run cmd/server/main.go"
        echo
        echo "  # Terminal 2: Run the test"
        echo "  ./scripts/test-graceful-shutdown.sh sigterm"
        exit 0
        ;;
    *)
        run_tests "$1"
        ;;
esac

print_success "All tests completed successfully! 🎉"
#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DB_NAME="oriondns_test"
TEST_DB_USER="postgres"
TEST_DB_HOST="localhost"
TEST_DB_PORT="5432"
TEST_DB_URL="postgres://${TEST_DB_USER}@${TEST_DB_HOST}:${TEST_DB_PORT}/${TEST_DB_NAME}?sslmode=disable"

print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

check_postgres() {
    print_header "Checking PostgreSQL Connection"
    
    if ! command -v psql &> /dev/null; then
        print_error "psql command not found. Please install PostgreSQL client."
        exit 1
    fi
    
    if ! pg_isready -h "$TEST_DB_HOST" -p "$TEST_DB_PORT" -U "$TEST_DB_USER" &> /dev/null; then
        print_error "PostgreSQL is not running or not accessible at ${TEST_DB_HOST}:${TEST_DB_PORT}"
        print_warning "Please start PostgreSQL and ensure user '$TEST_DB_USER' can connect"
        exit 1
    fi
    
    print_success "PostgreSQL is running and accessible"
}

setup_test_database() {
    print_header "Setting up Test Database"
    
    # Drop and recreate test database
    echo "Dropping existing test database (if exists)..."
    psql -h "$TEST_DB_HOST" -p "$TEST_DB_PORT" -U "$TEST_DB_USER" -c "DROP DATABASE IF EXISTS $TEST_DB_NAME;" postgres 2>/dev/null || true
    
    echo "Creating test database..."
    psql -h "$TEST_DB_HOST" -p "$TEST_DB_PORT" -U "$TEST_DB_USER" -c "CREATE DATABASE $TEST_DB_NAME;" postgres
    
    print_success "Test database '$TEST_DB_NAME' created"
}

run_migrations() {
    print_header "Running Database Migrations"
    
    if [ -f "migrations/tern-test.conf" ]; then
        cd migrations
        if command -v tern &> /dev/null; then
            tern migrate --config tern-test.conf
            print_success "Migrations completed"
        else
            print_warning "tern not found, skipping migrations"
            print_warning "You may need to run migrations manually"
        fi
        cd ..
    else
        print_warning "tern-test.conf not found, creating basic migration config"
        
        # Create basic tern config for testing
        cat > migrations/tern-test.conf << EOF
[database]
host = $TEST_DB_HOST
port = $TEST_DB_PORT
database = $TEST_DB_NAME
user = $TEST_DB_USER
sslmode = disable
EOF
        
        if command -v tern &> /dev/null; then
            cd migrations
            tern migrate --config tern-test.conf
            cd ..
            print_success "Migrations completed with generated config"
        else
            print_warning "tern not found, please install tern or run migrations manually"
        fi
    fi
}

run_unit_tests() {
    print_header "Running Unit Tests"
    
    echo "Running unit tests (database-independent)..."
    SKIP_DB_TESTS=true go test -v ./internal/... ./server/... 2>&1 | tee test_unit.log
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        print_success "Unit tests passed"
    else
        print_error "Unit tests failed"
        return 1
    fi
}

run_integration_tests() {
    print_header "Running Integration Tests"
    
    echo "Running integration tests (with database)..."
    TEST_DATABASE_URL="$TEST_DB_URL" go test -v ./... 2>&1 | tee test_integration.log
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        print_success "Integration tests passed"
    else
        print_error "Integration tests failed"
        return 1
    fi
}

run_coverage() {
    print_header "Running Coverage Analysis"
    
    echo "Generating test coverage report..."
    TEST_DATABASE_URL="$TEST_DB_URL" go test -v -coverprofile=coverage.out ./...
    
    if [ $? -eq 0 ]; then
        echo "Coverage by function:"
        go tool cover -func=coverage.out
        
        echo ""
        echo "Generating HTML coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        
        print_success "Coverage report generated: coverage.html"
        
        # Show total coverage
        TOTAL_COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}')
        echo -e "${GREEN}Total Coverage: $TOTAL_COVERAGE${NC}"
    else
        print_error "Coverage generation failed"
        return 1
    fi
}

cleanup() {
    print_header "Cleaning up"
    
    rm -f test_unit.log test_integration.log
    print_success "Cleanup completed"
}

show_help() {
    echo "OrionDNS Test Runner"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  unit          Run unit tests only (no database)"
    echo "  integration   Run integration tests (requires database)"
    echo "  coverage      Run tests with coverage analysis"
    echo "  setup         Setup test database and run migrations"
    echo "  all           Run complete test suite (default)"
    echo "  clean         Clean up test artifacts"
    echo "  help          Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  TEST_DB_URL   Custom test database URL"
    echo "  SKIP_DB_TESTS Set to 'true' to skip database tests"
}

main() {
    local command=${1:-all}
    
    case $command in
        "unit")
            run_unit_tests
            ;;
        "integration")
            check_postgres
            setup_test_database
            run_migrations
            run_integration_tests
            ;;
        "coverage")
            check_postgres
            setup_test_database
            run_migrations
            run_coverage
            ;;
        "setup")
            check_postgres
            setup_test_database
            run_migrations
            ;;
        "all")
            check_postgres
            setup_test_database
            run_migrations
            run_unit_tests
            run_integration_tests
            cleanup
            ;;
        "clean")
            cleanup
            rm -f coverage.out coverage.html
            print_success "All test artifacts cleaned"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
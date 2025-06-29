# Testing Guide for OrionDNS Backend

This guide provides comprehensive information about testing the OrionDNS backend, including setup, running tests, and contributing test code.

## Overview

The OrionDNS backend uses a multi-layered testing approach:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions with real database
- **DNS Integration Tests**: Test DNS server functionality
- **HTTP API Tests**: Test REST endpoints with mocked dependencies

## Test Structure

```
backend/
├── internal/
│   ├── testutil/          # Test utilities and fixtures
│   ├── */
│   │   └── *_test.go      # Unit tests for each module
├── server/
│   ├── dns/
│   │   ├── dns_test.go              # DNS unit tests
│   │   └── dns_integration_test.go  # DNS integration tests
│   └── web/
│       └── stats_test.go   # HTTP handler tests
├── integration_test.go     # Full integration tests
├── Makefile               # Test automation
└── test.sh               # Test runner script
```

## Prerequisites

### Required Software

- **Go 1.21+**: For running tests
- **PostgreSQL**: For integration tests
- **Make**: For using Makefile commands (optional)

### Database Setup

For integration tests, you need a PostgreSQL instance running:

```bash
# Using Docker (recommended)
docker run --name postgres-test -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15

# Or install PostgreSQL locally
# Ubuntu/Debian: sudo apt install postgresql postgresql-contrib
# macOS: brew install postgresql
```

## Running Tests

### Quick Start

```bash
# Run all tests (requires database)
make test

# Run only unit tests (no database required)
make test-unit

# Run with coverage report
make test-coverage
```

### Using the Test Script

The `test.sh` script provides a comprehensive test runner:

```bash
# Make script executable (one time)
chmod +x test.sh

# Run complete test suite
./test.sh all

# Run only unit tests
./test.sh unit

# Run only integration tests
./test.sh integration

# Generate coverage report
./test.sh coverage

# Setup test database only
./test.sh setup

# Clean test artifacts
./test.sh clean

# Show help
./test.sh help
```

### Manual Test Execution

```bash
# Unit tests only (no database)
SKIP_DB_TESTS=true go test -v ./internal/... ./server/...

# Integration tests (requires database)
TEST_DATABASE_URL="postgres://postgres@localhost:5432/oriondns_test?sslmode=disable" \
go test -v ./...

# Specific package tests
go test -v ./internal/stats/
go test -v ./server/dns/

# With race detection
go test -v -race ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TEST_DATABASE_URL` | PostgreSQL connection string for tests | `postgres://postgres@localhost:5432/oriondns_test?sslmode=disable` |
| `SKIP_DB_TESTS` | Skip database-dependent tests | `false` |

## Test Categories

### Unit Tests

Test individual functions and methods in isolation using mocks:

**Location**: `internal/*/db_test.go`, `server/dns/dns_test.go`, `server/web/stats_test.go`

**Characteristics**:
- Fast execution (< 1 second)
- No external dependencies
- Use mock objects for database/network calls
- Test specific business logic

**Example**:
```go
func TestStatsDB_Insert(t *testing.T) {
    // Setup mocks
    mockDB := &MockDB{}
    statsDB := stats.New(mockDB)
    
    // Test specific functionality
    err := statsDB.Insert(ctx, time.Now(), "google.com", "A")
    assert.NoError(t, err)
    
    // Verify mock expectations
    mockDB.AssertExpectations(t)
}
```

### Integration Tests

Test component interactions with real database:

**Location**: `internal/*/db_test.go`, `integration_test.go`

**Characteristics**:
- Require PostgreSQL database
- Test database operations end-to-end
- Use real database transactions
- Slower execution (1-10 seconds)

**Example**:
```go
func TestStatsDB_GetMostUsedDomains(t *testing.T) {
    testutil.SkipIfNoDatabase(t)
    pool := testutil.SetupTestDB(t)
    testutil.TruncateAllTables(t, pool)
    
    // Insert test data
    database := db.New(pool)
    statsDB := stats.New(database)
    
    // Test database operations
    results, err := statsDB.GetMostUsedDomains(ctx, from, to, categories, 10)
    require.NoError(t, err)
    assert.Len(t, results, expectedCount)
}
```

### DNS Integration Tests

Test DNS server logic with various blocking scenarios:

**Location**: `server/dns/dns_integration_test.go`

**Characteristics**:
- Test DNS message processing
- Test domain blocking logic
- Test caching behavior
- Test concurrent access patterns

### HTTP API Tests

Test REST endpoints with mocked dependencies:

**Location**: `server/web/*_test.go`

**Characteristics**:
- Test HTTP request/response handling
- Test JSON serialization/deserialization
- Test error handling
- Use mock database interfaces

## Test Utilities

### `testutil` Package

The `internal/testutil` package provides common testing utilities:

**Database Setup**:
```go
func TestExample(t *testing.T) {
    testutil.SkipIfNoDatabase(t)           // Skip if no DB available
    pool := testutil.SetupTestDB(t)        // Setup test database
    testutil.TruncateAllTables(t, pool)    // Clean tables
    
    // Your test code here
}
```

**Test Fixtures**:
```go
// Create test data
blockedDomains := testutil.CreateTestBlockedDomains()
categories := testutil.CreateTestCategories()
domains := testutil.CreateTestDomains()
```

### Mock Objects

Create mocks for interfaces using testify/mock:

```go
type MockStatsDB struct {
    mock.Mock
}

func (m *MockStatsDB) Insert(ctx context.Context, t time.Time, domain, domainType string) error {
    args := m.Called(ctx, t, domain, domainType)
    return args.Error(0)
}

// In test:
mockStats := &MockStatsDB{}
mockStats.On("Insert", mock.Anything, mock.Anything, "google.com", "A").Return(nil)
```

## Writing New Tests

### Test Naming Conventions

- Test files: `*_test.go`
- Test functions: `TestFunction_Scenario(t *testing.T)`
- Benchmark functions: `BenchmarkFunction_Scenario(b *testing.B)`

### Test Structure

Follow the Arrange-Act-Assert pattern:

```go
func TestFunction_Scenario(t *testing.T) {
    // Arrange - Setup test data and dependencies
    testData := createTestData()
    mockDep := &MockDependency{}
    
    // Act - Execute the function under test
    result, err := functionUnderTest(testData)
    
    // Assert - Verify the results
    require.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}
```

### Database Tests

For tests requiring database:

1. **Always** call `testutil.SkipIfNoDatabase(t)` first
2. **Always** setup and cleanup database state
3. **Always** use `testutil.TruncateAllTables(t, pool)` for isolation
4. Use `require.NoError(t, err)` for setup operations
5. Use `assert.*` for test assertions

### DNS Tests

For DNS-related tests:

1. Create DNS messages using `github.com/miekg/dns`
2. Test both exact and recursive domain matching
3. Test cache behavior
4. Test concurrent access patterns

### HTTP Tests

For HTTP endpoint tests:

1. Use `httptest.NewRequest()` and `httptest.NewRecorder()`
2. Mock all database dependencies
3. Test both success and error cases
4. Verify HTTP status codes and response bodies

## Coverage Requirements

### Target Coverage

- **Overall**: 80%+
- **Critical components** (DNS logic, database layer): 90%+
- **HTTP handlers**: 85%+

### Generating Coverage Reports

```bash
# Generate coverage report
make test-coverage

# View in browser (opens coverage.html)
open coverage.html

# Terminal output
make test-coverage-text
```

### Coverage Guidelines

- Focus on testing business logic over boilerplate
- Ensure all error paths are tested
- Test edge cases and boundary conditions
- Don't aim for 100% coverage - focus on critical paths

## Continuous Integration

### GitHub Actions

The project includes a comprehensive CI pipeline in `.github/workflows/test.yml`:

**Jobs**:
- **Unit Tests**: Run on Go 1.21, 1.22, 1.23
- **Integration Tests**: Run with PostgreSQL service
- **Coverage**: Generate and upload coverage reports
- **Lint**: Run golangci-lint
- **Security**: Run Gosec security scanner

**Triggers**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Changes in `backend/` directory

### Local CI Simulation

```bash
# Run the same checks as CI
make ci-test

# Or individually:
make lint        # Linting
make vet         # Go vet
make test        # All tests
make test-coverage  # Coverage
```

## Troubleshooting

### Common Issues

**Database Connection Errors**:
```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432 -U postgres

# Create test database manually
psql -h localhost -U postgres -c "CREATE DATABASE oriondns_test;"
```

**Test Timeouts**:
```bash
# Increase timeout for slow tests
go test -timeout 30s ./...
```

**Race Condition Detection**:
```bash
# Run with race detector
go test -race ./...
```

### Test Data Isolation

- Each test should be independent
- Use `testutil.TruncateAllTables()` between tests
- Don't rely on test execution order
- Clean up resources in test cleanup functions

### Performance Considerations

- Unit tests should complete in < 100ms
- Integration tests should complete in < 5s
- Use `t.Parallel()` for independent tests
- Profile slow tests with `go test -bench=.`

## Best Practices

1. **Test Independence**: Each test should be able to run independently
2. **Clear Test Names**: Use descriptive test function names
3. **Arrange-Act-Assert**: Follow this pattern consistently
4. **Error Testing**: Test both success and failure scenarios
5. **Mock External Dependencies**: Use mocks for external services
6. **Table-Driven Tests**: Use for testing multiple scenarios
7. **Test Data**: Use the `testutil` package for consistent test data
8. **Documentation**: Comment complex test logic

## Contributing Tests

When contributing new tests:

1. Follow existing patterns and conventions
2. Ensure tests pass locally before submitting
3. Include both positive and negative test cases
4. Update this documentation if adding new testing patterns
5. Ensure CI pipeline passes

For questions about testing, refer to the main project documentation or open an issue.
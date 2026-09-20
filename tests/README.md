# Tests

This directory contains tests for the Docker Azure Key Vault Secrets Provider.

## Test Organization

### Unit Tests (in package directories)

**`internal/config/config_test.go`**
- Tests for configuration loading from environment variables
- Validates error handling for missing or invalid configuration
- Tests various URL formats

**`internal/provider/provider_test.go`**
- Tests for secret name parsing (`parseSecretName`)
- Validates input validation and error cases
- Tests provider creation and lifecycle

### Integration Tests

**`tests/integration_test.go`**
- End-to-end workflow tests
- Demonstrates typical provider usage patterns
- Performance benchmarks

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Specific Package Tests

```bash
# Config package tests
go test ./internal/config

# Provider package tests
go test ./internal/provider

# Integration tests
go test ./tests
```

### Run Specific Test

```bash
go test -run TestParseSecretName ./internal/provider
```

### Verbose Output

```bash
go test -v ./...
```

### With Coverage

```bash
go test -cover ./...
```

### Generate Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View as HTML
go tool cover -html=coverage.out

# View text summary
go tool cover -func=coverage.out
```

### Run Benchmarks

```bash
go test -bench=. ./tests
```

### Run Benchmarks with Memory Stats

```bash
go test -bench=. -benchmem ./tests
```

## Test Coverage

Current test coverage includes:

| Package | Tests | Coverage |
|---------|-------|----------|
| `internal/config` | 5 test functions | Secret name validation, configuration loading |
| `internal/provider` | 3 test functions | Provider lifecycle, secret parsing |
| `tests` | 3 test functions | Integration workflows, benchmarks |
| **Total** | **11+ tests** | **Core functionality** |

## Test Cases

### Config Tests
- ✓ Valid configuration loading
- ✓ Missing environment variable handling
- ✓ Empty environment variable handling
- ✓ Multiple URL formats
- ✓ Repeated configuration loading

### Provider Tests
- ✓ Valid secret name parsing (12 cases)
- ✓ Invalid prefix handling
- ✓ Empty secret name handling
- ✓ Nested path rejection
- ✓ Provider creation
- ✓ Provider lifecycle (Run/shutdown)

### Integration Tests
- ✓ Basic provider workflow
- ✓ Configuration integration
- ✓ Performance benchmarks

## Adding New Tests

When adding new functionality:

1. Create test file in the same package: `filename_test.go`
2. Use table-driven tests for multiple scenarios
3. Follow naming convention: `Test<FunctionName>`
4. Add comments explaining the test purpose
5. Run full test suite to ensure no regressions

### Example Test

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid case",
            input:   "test",
            want:    "result",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Continuous Integration

When pushing code:

```bash
# Run all tests before committing
go test ./...

# Check coverage
go test -cover ./...

# Lint code
go vet ./...
go fmt ./...
```

## Mocking and Testing with Azure SDK

For tests that require Azure SDK clients, you can:

1. Use `nil` clients for validation tests
2. Create mock implementations of client interfaces
3. Use dependency injection to swap implementations

Example:

```go
type MockSecretsClient struct {
    GetSecretFunc func(ctx context.Context, name string) (*SecretResponse, error)
}

func (m *MockSecretsClient) GetSecret(ctx context.Context, name string, ...) (*SecretResponse, error) {
    return m.GetSecretFunc(ctx, name)
}
```

## Performance Benchmarks

Run benchmarks to measure performance:

```bash
go test -bench=BenchmarkProviderCreation -benchmem ./tests
```

Expected output:
```
BenchmarkProviderCreation-8    1000000     1234 ns/op     0 B/op     0 allocs/op
```

## Troubleshooting Tests

### Test Fails with "AZURE_KEYVAULT_URL is required"

This is expected when environment variables aren't set. Tests are designed to work in any environment.

### Test Timeout

If a test hangs:

```bash
# Run with timeout (default 10m)
go test -timeout 30s ./...
```

### Test Data Issues

Ensure tests don't depend on external state:
- Use `defer` to clean up environment variables
- Each test should be independent
- Use `t.Parallel()` for concurrent tests

## Future Test Improvements

- [ ] Add GetSecrets mock tests with Azure SDK responses
- [ ] Add performance regression tests
- [ ] Add fuzzing tests for input validation
- [ ] Add end-to-end tests with Docker
- [ ] Add stress tests for concurrent secret retrieval

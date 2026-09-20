# Contributing

Thanks for your interest in contributing to the Docker Azure Key Vault Secrets Provider! This guide will help you get started.

## Code of Conduct

Be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- Go 1.26 or later
- Git
- Azure subscription (for testing)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/docker-azure-keyvault-provider.git
cd docker-azure-keyvault-provider

# Install dependencies
go mod download

# Build
go build -o provider ./cmd/provider

# Run tests
go test ./...
```

## Project Structure

```
.
├── cmd/
│   └── provider/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration and environment loading
│   ├── keyvault/
│   │   └── azure_kv_client.go # Azure SDK integration (future)
│   └── provider/
│       └── provider.go       # Core provider logic
├── README.md                 # Main documentation
├── QUICKSTART.md             # Quick start guide
├── TROUBLESHOOTING.md        # Troubleshooting guide
└── go.mod                    # Go module file
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

Follow these guidelines:

#### Code Style
- Use `gofmt` for formatting: `go fmt ./...`
- Use `golint` for linting: `go vet ./...`
- Write clear, descriptive variable and function names
- Keep functions small and focused

#### Documentation
- Add GoDoc comments to all public functions and types
- Update README.md if behavior changes
- Add comments for complex logic

#### Testing
- Write unit tests for new functionality
- Ensure all existing tests pass: `go test ./...`
- Aim for >80% code coverage for new code

#### Example: Adding a New Feature

```go
// config.go
// Add GoDoc comment
// CacheTTL is the time-to-live for cached secrets
const CacheTTL = 5 * time.Minute

// provider.go
// Add GoDoc comment
// WithCache enables secret caching with the given TTL
func (p *AzureKeyVaultProvider) WithCache(ttl time.Duration) {
    // implementation
}
```

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/provider

# Run a specific test
go test -run TestParseSecretName ./internal/provider
```

### 4. Format and Lint

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Install golangci-lint for comprehensive linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

### 5. Commit Changes

```bash
git add .
git commit -m "feat: add secret caching support"
# or
git commit -m "fix: handle empty secret names correctly"
```

#### Commit Message Format

Use conventional commits:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, missing semicolons, etc.)
- `refactor:` - Code refactoring without feature changes
- `test:` - Adding or updating tests
- `chore:` - Build, dependency, or CI changes

Example:
```
feat: add support for secret versioning

- Allow retrieving specific secret versions
- Add version parameter to GetSecrets
- Update documentation
```

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub:
- Describe what changes you made
- Reference any related issues
- Ensure CI checks pass

## Common Tasks

### Running Tests in Watch Mode

```bash
# Install entr
go install github.com/cortesi/entr/cmd/entr@latest

# Watch for changes and run tests
ls internal/**/*.go | entr -c go test ./...
```

### Debugging

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o provider ./cmd/provider

# Run with debugger (dlv)
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/provider
```

### Generating Test Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View as HTML
go tool cover -html=coverage.out

# View text summary
go tool cover -func=coverage.out
```

## Areas for Contribution

### High Priority
- [ ] Add structured logging support
- [ ] Implement request timeouts
- [ ] Add unit tests for all functions
- [ ] Support for secret versioning
- [ ] Add health check endpoint

### Medium Priority
- [ ] Implement secret caching
- [ ] Add metrics/instrumentation
- [ ] Support for managed identity scenarios
- [ ] Rate limiting
- [ ] Better error messages

### Low Priority
- [ ] Performance optimizations
- [ ] Additional documentation
- [ ] Example Docker Compose files
- [ ] Integration tests

## Pull Request Guidelines

### Before Submitting

- [ ] Code is formatted (`go fmt ./...`)
- [ ] Code passes linting (`go vet ./...`)
- [ ] All tests pass (`go test ./...`)
- [ ] New tests are added for new functionality
- [ ] Documentation is updated
- [ ] Commit messages follow the format
- [ ] No merge conflicts

### Review Process

1. Automated checks must pass (format, lint, tests)
2. At least one maintainer review
3. Address review feedback
4. Approved and merged

### After Merge

- Delete your feature branch
- The code will be included in the next release

## Release Process

Releases follow semantic versioning (MAJOR.MINOR.PATCH):

- MAJOR: Breaking changes
- MINOR: New features
- PATCH: Bug fixes

Example: v1.2.3

## Questions?

- Check [README.md](README.md) for general information
- See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues
- Open an issue for questions or problems
- Reach out to maintainers in discussions

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

Thank you for contributing! 🎉

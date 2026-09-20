# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial version with basic secret retrieval functionality
- Support for Azure Key Vault integration
- Docker Secrets Engine plugin interface implementation
- GoDoc documentation for public API
- Configuration via environment variables
- Comprehensive README, QUICKSTART, and TROUBLESHOOTING guides
- Contributing guidelines

### Changed
- N/A

### Fixed
- N/A

### Security
- N/A

### Deprecated
- N/A

### Removed
- N/A

---

## Planned Features

### v0.2.0
- [ ] Structured logging support
- [ ] Request timeout handling
- [ ] Comprehensive unit tests
- [ ] Secret versioning support
- [ ] Health check endpoint

### v0.3.0
- [ ] Secret caching with TTL
- [ ] Metrics and instrumentation
- [ ] Rate limiting
- [ ] Graceful shutdown

### v1.0.0
- Production-ready release
- All security best practices implemented
- Complete test coverage
- Performance optimizations

---

## Version History

### How to Document Changes

When releasing a new version:

1. Move items from "Unreleased" to a new section with the version and date
2. Follow these categories:
   - **Added** for new features
   - **Changed** for changes in existing functionality
   - **Deprecated** for soon-to-be removed features
   - **Removed** for now removed features
   - **Fixed** for any bug fixes
   - **Security** for security vulnerability fixes

3. Format: `## [version] - YYYY-MM-DD`

### Example Entry

```markdown
## [0.1.0] - 2026-09-19

### Added
- Initial release with basic secret retrieval
- Docker Secrets Engine plugin support
- Azure Key Vault integration

### Fixed
- Handle empty secret values gracefully
```

---

For releases, also update version in:
- `cmd/provider/main.go` (if version constant is added)
- `go.mod` (if applicable)
- Docker image tags (if containerized)

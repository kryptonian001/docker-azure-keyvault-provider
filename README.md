# Docker Azure Key Vault Secrets Provider

A Docker Secrets Engine plugin that retrieves secrets from Azure Key Vault.

## Quick Links

📖 **[Full Documentation](docs/README.md)** — Complete guide and API reference  
🚀 **[Quick Start](docs/QUICKSTART.md)** — Get running in 5 minutes  
❓ **[Troubleshooting](docs/TROUBLESHOOTING.md)** — Common issues and solutions  
🤝 **[Contributing](docs/CONTRIBUTING.md)** — Developer guidelines  
📝 **[Changelog](docs/CHANGELOG.md)** — Version history  
⚖️ **[License](docs/LICENSE)** — License information  

## Overview

This plugin enables Docker to use secrets stored in Azure Key Vault. When Docker needs a secret, the provider retrieves it from your Key Vault using the secret ID format: `azure/<secret-name>`

## Key Features

- 🔐 Azure Key Vault integration for Docker Secrets Engine
- 🔑 Support for managed identity and service principal authentication
- 📦 Simple secret ID format: `azure/<secret-name>`
- 🚀 Production-ready with comprehensive documentation
- 📚 Full API documentation with GoDoc comments

## Getting Started

### 1. Prerequisites

- Windows 10+ or macOS/Linux with Docker Desktop
- Azure Key Vault instance
- Azure credentials (managed identity or service principal)

### 2. Quick Setup

Set up Azure credentials:
```powershell
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
$env:AZURE_CLIENT_ID="<your-client-id>"
$env:AZURE_CLIENT_SECRET="<your-client-secret>"
$env:AZURE_TENANT_ID="<your-tenant-id>"
```

Build and run:
```bash
go build -o provider ./cmd/provider
./provider
```

### 3. Use in Docker

```yaml
version: '3.8'
services:
  app:
    image: myapp:latest
    secrets:
      - db_password

secrets:
  db_password:
    external: true
    name: azure/my-database-password
```

## Documentation

All documentation has been organized in the `docs/` folder:

| Document | Purpose |
|----------|---------|
| [docs/README.md](docs/README.md) | Complete reference guide |
| [docs/QUICKSTART.md](docs/QUICKSTART.md) | 5-minute setup guide |
| [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) | Issue resolution guide |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | Development guidelines |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Release notes |

## Project Structure

```
.
├── cmd/
│   └── provider/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loading
│   ├── keyvault/
│   │   └── azure_kv_client.go # Azure SDK integration
│   └── provider/
│       └── provider.go       # Core provider logic
├── docs/                     # Documentation
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── TROUBLESHOOTING.md
│   ├── CONTRIBUTING.md
│   ├── CHANGELOG.md
│   └── LICENSE
├── go.mod
└── go.sum
```

## Security

- Never commit credentials to version control
- Use managed identity in production
- Follow Azure Key Vault access best practices
- Rotate credentials regularly

For more details, see [docs/README.md](docs/README.md#security-considerations).

## Support

- Check [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for common issues
- Review [docs/README.md](docs/README.md) for comprehensive documentation
- See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for development information

## License

See [docs/LICENSE](docs/LICENSE) for license information.

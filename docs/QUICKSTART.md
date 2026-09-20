# Quick Start Guide

Get the Azure Key Vault Secrets Provider running in 5 minutes.

## 1. Prerequisites

- Windows 10+ or macOS/Linux with Docker Desktop
- Azure subscription with Key Vault
- Go 1.26+ (for building)

## 2. Prepare Azure Credentials

Choose one method:

### Method A: Service Principal (Development)
```powershell
# Set environment variables
$env:AZURE_CLIENT_ID="<your-client-id>"
$env:AZURE_CLIENT_SECRET="<your-client-secret>"
$env:AZURE_TENANT_ID="<your-tenant-id>"
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
```

### Method B: Managed Identity (Production on Azure)
```powershell
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
```

### Method C: Azure CLI (Local Development)
```powershell
az login
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
```

## 3. Create a Test Secret in Azure Key Vault

```bash
# Using Azure CLI
az keyvault secret set --vault-name myvault --name test-secret --value "my-secret-value"
```

## 4. Build the Provider

```bash
cd docker-azure-keyvault-provider
go build -o provider.exe ./cmd/provider
```

## 5. Run the Provider

The provider connects to Docker Secrets Engine. On Docker Desktop, this happens automatically:

```bash
.\provider.exe
```

You should see:
```
2026/09/19 19:34:51 starting Azure Key Vault Secrets Engine provider
```

## 6. Use Secrets in Docker

### Create Docker Compose File

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  app:
    image: alpine:latest
    command: sh -c "cat /run/secrets/test-secret"
    secrets:
      - test-secret

secrets:
  test-secret:
    external: true
    name: azure/test-secret
```

### Run Docker Compose

In another terminal:
```bash
docker compose up
```

You should see the secret value printed:
```
my-secret-value
```

## 7. Verify It Works

```bash
# Check if provider is running
# Look for the provider process in Task Manager or:
# docker ps (provider may be shown as a plugin)

# Test retrieving a secret
# The secret should be accessible in your Docker containers
```

## Next Steps

- Read [README.md](README.md) for detailed configuration options
- Check [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
- See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues

## Common Issues

### Provider won't start: "AZURE_KEYVAULT_URL is required"
→ Set the environment variable: `$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"`

### Secret not found: "failed to retrieve secret from Azure Key Vault"
→ Verify the secret exists in your Key Vault and your credentials have `Get` permission

### Connection refused: "failed to connect to Docker Secrets Engine"
→ Ensure Docker Desktop is running and Secrets Engine is enabled

For more help, see [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

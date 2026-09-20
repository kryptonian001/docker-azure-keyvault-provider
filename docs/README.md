# Docker Azure Key Vault Secrets Provider

A Docker Secrets Engine plugin that retrieves secrets from Azure Key Vault.

## Overview

This plugin enables Docker to use secrets stored in Azure Key Vault. When Docker needs a secret, the provider retrieves it from Azure Key Vault using the secret ID.

**Secret ID Format:** `azure/<secret-name>`

Example: `azure/my-database-password` retrieves the secret named `my-database-password` from Azure Key Vault.

## Prerequisites

- Docker Desktop with Secrets Engine support
- Azure account with Key Vault instance
- Azure credentials (managed identity or service principal)
- Go 1.26+ (for building from source)

## Installation

### 1. Set Up Azure Credentials

The plugin uses Azure's default credential chain. Set up one of the following:

**Option A: Environment Variables (Service Principal)**
```bash
# Windows PowerShell
$env:AZURE_CLIENT_ID="<your-client-id>"
$env:AZURE_CLIENT_SECRET="<your-client-secret>"
$env:AZURE_TENANT_ID="<your-tenant-id>"
$env:AZURE_KEYVAULT_URL="https://<vault-name>.vault.azure.net/"
```

**Option B: Managed Identity (Azure VM/Container)**
```bash
$env:AZURE_KEYVAULT_URL="https://<vault-name>.vault.azure.net/"
```

**Option C: Azure CLI (local development)**
```bash
az login
$env:AZURE_KEYVAULT_URL="https://<vault-name>.vault.azure.net/"
```

**Option D: Device Code Login (Headless/Remote)**
```bash
# Use device code authentication for headless environments
az login --use-device-code
$env:AZURE_KEYVAULT_URL="https://<vault-name>.vault.azure.net/"
```
Device code login is ideal for:
- Headless or remote environments without browser access
- SSH sessions or restricted network environments
- CI/CD pipelines and automated deployments
- Scenarios where interactive authentication isn't possible

### 2. Build the Plugin

```bash
go build -o provider.exe ./cmd/provider
```

### 3. Run the Plugin

The plugin must be run with Docker Secrets Engine. Docker Desktop manages the connection automatically.

```bash
.\provider.exe
```

## Configuration

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `AZURE_KEYVAULT_URL` | URL to your Azure Key Vault | `https://myvault.vault.azure.net/` |

### Azure Authentication

The plugin uses Azure SDK's default credential chain (in order):
1. Environment variables (`AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`)
2. Managed Identity (if running on Azure resources)
3. Azure CLI credentials (including device code login)
4. Visual Studio credentials

#### Device Code Login
Device code login provides an alternative authentication method for scenarios where traditional interactive login isn't available:

```bash
# Authenticate using device code
az login --use-device-code

# A device code will be displayed - open the URL in another device's browser
# and enter the code when prompted
```

Once authenticated, the provider automatically uses the cached Azure CLI credentials.

## Usage

### Using Secrets with Docker Compose

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

### Using Secrets in Docker Stack

```bash
docker secret create db_password --external-provider azure/my-database-password
```

### Accessing Secrets in Containers

Secrets are mounted at `/run/secrets/<secret-name>`:

```bash
# Inside container
cat /run/secrets/db_password
```

## Secret Naming Rules

- **Required Prefix:** All secret IDs must start with `azure/`
- **No Nested Paths:** Slashes are not allowed in the secret name (e.g., `azure/folder/secret` will fail)
- **Valid Examples:**
  - `azure/database-password` ✓
  - `azure/api-key-prod` ✓
  - `azure/db_connection_string` ✓
- **Invalid Examples:**
  - `azure/folder/secret` ✗ (nested paths not supported)
  - `mysql-password` ✗ (missing `azure/` prefix)

## Development

### Project Structure

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
│       └── provider.go       # Provider implementation
└── go.mod                    # Go module definition
```

### Building

```bash
go build -o provider.exe ./cmd/provider
```

### Testing

```bash
go test ./...
```

### Environment Setup for Local Development

```powershell
# Install Go (if not already installed)
# https://golang.org/dl

# Clone repository
git clone <repository-url>
cd docker-azure-keyvault-provider

# Set Azure credentials for testing
$env:AZURE_KEYVAULT_URL="https://your-vault.vault.azure.net/"
$env:AZURE_CLIENT_ID="your-client-id"
$env:AZURE_CLIENT_SECRET="your-client-secret"
$env:AZURE_TENANT_ID="your-tenant-id"

# Build
go build ./cmd/provider
```

## Troubleshooting

### Error: "AZURE_KEYVAULT_URL is required"

**Cause:** Environment variable is not set  
**Fix:** Set the `AZURE_KEYVAULT_URL` environment variable

```powershell
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
```

### Error: "failed to create Azure credentials"

**Cause:** Azure credentials are not configured or invalid  
**Fix:** 
1. Verify credentials are set correctly
2. Check Azure CLI login: `az login`
3. Ensure account has access to the Key Vault

### Error: "failed to retrieve secret from Azure Key Vault"

**Cause:** Secret doesn't exist or access is denied  
**Fix:**
1. Verify secret name in Key Vault (without `azure/` prefix)
2. Check access permissions: `az keyvault secret show --vault-name <vault> --name <secret>`
3. Ensure Key Vault access policy grants `Get` permission

### Error: "unsupported secret ID: expected prefix azure/"

**Cause:** Secret ID doesn't have the required `azure/` prefix  
**Fix:** Use the correct format: `azure/<secret-name>`

## Security Considerations

1. **Credentials:** Never commit credentials to version control. Use environment variables or managed identities.
2. **Key Vault Access:** Follow Azure Key Vault access best practices:
   - Use RBAC or access policies to limit permissions
   - Grant only `Get` permission for required secrets
3. **Secret Rotation:** Rotate credentials regularly in Azure Key Vault
4. **Logging:** The plugin logs to stdout. Be careful not to expose secrets in logs.

## License

[Add your license here]

## Support

For issues, questions, or contributions, please refer to the project repository.

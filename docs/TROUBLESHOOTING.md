# Troubleshooting Guide

## Configuration Issues

### Error: "AZURE_KEYVAULT_URL is required"

**Cause:** Environment variable not set

**Solutions:**
```powershell
# Windows PowerShell
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"

# Windows Command Prompt
set AZURE_KEYVAULT_URL=https://myvault.vault.azure.net/

# macOS/Linux
export AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
```

**Verification:**
```powershell
# Check if set
$env:AZURE_KEYVAULT_URL
echo $AZURE_KEYVAULT_URL  # macOS/Linux
```

---

## Authentication Issues

### Error: "failed to create Azure credentials"

**Cause:** Azure credentials not found or invalid

**Check your auth method:**

```powershell
# Method 1: Service Principal (check variables)
$env:AZURE_CLIENT_ID
$env:AZURE_CLIENT_SECRET
$env:AZURE_TENANT_ID

# Method 2: Azure CLI login
az account show  # Should show your account

# Method 3: Check managed identity (on Azure VMs)
Invoke-WebRequest -Uri "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2017-09-01&resource=https://vault.azure.net" -Headers @{Metadata="true"}
```

**Solutions:**

**Option A: Use Service Principal**
```powershell
# Create a service principal
az ad sp create-for-rbac --name "docker-azure-kv-provider" --role Contributor

# Set environment variables
$env:AZURE_CLIENT_ID="<appId>"
$env:AZURE_CLIENT_SECRET="<password>"
$env:AZURE_TENANT_ID="<tenant>"
```

**Option B: Use Azure CLI**
```powershell
az login
# Credentials are cached, provider will use them automatically
```

**Option C: Use Managed Identity (Azure VMs only)**
```powershell
# Ensure VM has managed identity assigned
# No environment variables needed
$env:AZURE_KEYVAULT_URL="https://myvault.vault.azure.net/"
```

---

## Key Vault Access Issues

### Error: "failed to retrieve secret from Azure Key Vault"

**Causes:**
1. Secret doesn't exist
2. No permission to access secret
3. Key Vault access is restricted

**Solutions:**

**1. Verify secret exists:**
```bash
az keyvault secret show --vault-name myvault --name test-secret
```

**2. Check permissions:**
```bash
# List your permissions
az keyvault secret list --vault-name myvault

# Grant permission to service principal
az keyvault set-policy --name myvault \
  --spn <client-id> \
  --secret-permissions get list
```

**3. Check Key Vault firewall:**
```bash
az keyvault show --name myvault --query "properties.networkAcls"

# If firewall is enabled, add your IP
az keyvault network-rule add --name myvault --ip-address <your-ip>
```

**4. Check Key Vault access:**
```bash
# Try accessing directly with Azure CLI
az keyvault secret show --vault-name myvault --name test-secret

# If this fails, it's an Azure permission issue, not the provider
```

---

## Secret ID Issues

### Error: "unsupported secret ID: expected prefix azure/"

**Cause:** Secret ID missing `azure/` prefix

**Solution:**
```yaml
# ❌ Wrong
secrets:
  my-secret:
    external: true
    name: test-secret

# ✓ Correct
secrets:
  my-secret:
    external: true
    name: azure/test-secret
```

### Error: "invalid secret name: nested paths are not supported"

**Cause:** Secret name contains `/`

**Examples:**
```
❌ azure/folder/secret  (nested path)
❌ azure/path/to/secret
✓ azure/mysecret
✓ azure/my-secret
✓ azure/my_secret_123
```

**Solution:** Use a single secret name without slashes

```bash
# Create secret without nesting
az keyvault secret set --vault-name myvault --name my-secret --value "value"

# Reference it correctly
azure/my-secret  # ✓ Correct
```

### Error: "Azure Key Vault returned an empty value for secret"

**Cause:** Secret exists but has no value

**Solution:**
```bash
# Delete and recreate the secret
az keyvault secret delete --vault-name myvault --name test-secret
az keyvault secret set --vault-name myvault --name test-secret --value "your-secret-value"
```

---

## Docker Connection Issues

### Error: "failed to connect to Docker Secrets Engine"

**Causes:**
1. Docker Desktop not running
2. Secrets Engine not enabled
3. Wrong socket path (Windows vs macOS/Linux)

**Solutions:**

**1. Ensure Docker Desktop is running:**
```powershell
docker --version  # Should show version
docker ps         # Should work
```

**2. Check Secrets Engine is enabled:**
- Windows/Mac: Docker Desktop Settings → Experimental Features → Enable Beta features

**3. Verify socket path:**
```powershell
# Windows (provider uses Docker socket automatically)
# macOS/Linux
ls /run/docker.sock  # Should exist
```

---

## Connection Issues

### Error: "connection refused" or "connection timeout"

**Causes:**
1. Docker daemon not running
2. Docker socket permissions
3. Network connectivity to Azure

**Solutions:**

**1. Restart Docker:**
```powershell
# Windows/Mac
# Restart Docker Desktop from taskbar/menu bar

# Linux
sudo systemctl restart docker
```

**2. Check Docker socket permissions (Linux):**
```bash
# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

**3. Test Azure connectivity:**
```bash
# Test connection to Key Vault
curl https://myvault.vault.azure.net/
# Should respond (may be 401 without auth, that's OK)
```

---

## Performance Issues

### Secrets retrieval is slow

**Causes:**
1. Network latency to Azure
2. Azure SDK initialization overhead

**Solutions:**

1. **Verify Azure connectivity:**
   ```bash
   ping vault.azure.net
   # Should have reasonable latency (< 100ms)
   ```

2. **Check provider logs:**
   - Provider logs to stdout
   - Capture output for timing information

3. **Network optimization:**
   - Use Azure region nearest to your location
   - Check VPN/firewall configuration
   - Verify ISP connectivity to Azure

---

## Debugging

### Enable verbose logging

The provider logs to stdout. Capture output:

```powershell
# Windows PowerShell
.\provider.exe | Tee-Object -FilePath provider.log

# Linux/macOS
./provider 2>&1 | tee provider.log
```

### Test Azure credentials directly

```powershell
# Verify Azure CLI auth works
az keyvault secret show --vault-name myvault --name test-secret

# If Azure CLI works but provider doesn't, it's a provider bug
# Report with the provider output
```

### Test with a simple secret

```bash
# Create a test secret with simple value
az keyvault secret set --vault-name myvault --name test-secret --value "hello"

# Try to use it in Docker
# If this works, other secrets should work too
```

---

## Still Having Issues?

1. Verify all environment variables are set correctly
2. Run the tests: `go test ./...`
3. Check the [README.md](README.md) for full documentation
4. Review Docker Secrets Engine documentation
5. Check Azure Key Vault access policies and firewall settings

**When reporting issues, include:**
- Provider version and build date
- Operating system
- Docker Desktop version
- Azure Key Vault region
- Relevant error messages and logs

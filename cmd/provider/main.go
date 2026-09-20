// Package main provides the Docker Secrets Engine plugin for Azure Key Vault.
//
// This program acts as a secrets provider for Docker, retrieving secrets from
// Azure Key Vault. It must be configured with AZURE_KEYVAULT_URL environment variable
// and appropriate Azure credentials (via environment variables or managed identity).
//
// Usage:
//
//	provider [flags]
//
// Environment Variables:
//
//	AZURE_KEYVAULT_URL - Required. URL to the Azure Key Vault instance
//	AZURE_CLIENT_ID - Azure Service Principal client ID (optional)
//	AZURE_CLIENT_SECRET - Azure Service Principal client secret (optional)
//	AZURE_TENANT_ID - Azure tenant ID (optional)
package main

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/secrets-engine/plugin"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/config"
	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/provider"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Create credentials for Azure
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to create Azure credentials: %v", err)
	}

	// Create the Azure Key Vault secrets client
	client, err := azsecrets.NewClient(cfg.VaultURL, credential, nil)
	if err != nil {
		log.Fatalf("failed to create Azure Key Vault client: %v", err)
	}

	// Create our Secrets Engine provider
	azureProvider := provider.New(client)

	// Docker Desktop exposes its Secrets Engine through engine.sock.
	socketPath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		"docker-secrets-engine",
		"engine.sock",
	)

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf(
			"failed to connect to Docker Secrets Engine at %s: %v",
			socketPath,
			err,
		)
	}
	defer conn.Close()

	// Configure our provider registration.
	pluginConfig := plugin.Config{
		Version: plugin.MustNewVersion("v0.0.1"),
		SecretsProviderConfig: &plugin.SecretsProviderConfig{
			Pattern: plugin.MustParsePattern("azure/**"),
		},
	}

	// Register the provider with Docker Secrets Engine.
	p, err := plugin.NewSecretsProvider(
		azureProvider,
		pluginConfig,
		plugin.WithConnection(conn),
		plugin.WithPluginName("azure-keyvault"),
	)
	if err != nil {
		log.Fatalf("failed to create plugin: %v", err)
	}

	log.Printf("starting Azure Key Vault Secrets Engine provider")

	// Keep the provider running and servicing secret requests.
	if err := p.Run(ctx); err != nil {
		log.Fatalf("plugin failed: %v", err)
	}
}

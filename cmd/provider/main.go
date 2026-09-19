package main

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/docker/secrets-engine/plugin"
	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/config"
	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/provider"
)

func main() {
	ctx := context.Background()

	vaultConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	credentials, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain credentials: %v", err)
	}

	client, err := azsecrets.NewClient(
		vaultConfig.VaultURL,
		credentials,
		nil,
	)
	if err != nil {
		log.Fatalf("Key Vault client creation failed: %v", err)
	}

	azureProvider :=
		provider.New(client)

	pluginConfig := plugin.Config{
		Version: plugin.MustNewVersion("v0.1.0"),

		SecretsProviderConfig: &plugin.SecretsProviderConfig{
			Pattern: plugin.MustParsePattern("azure/**"),
		},
	}

	p, err := plugin.NewSecretsProvider(
		azureProvider,
		pluginConfig,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	_ = ctx
	_ = client
}

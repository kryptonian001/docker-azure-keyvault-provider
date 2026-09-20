// Package config provides configuration loading for the Azure Key Vault provider.
package config

import (
	"fmt"
	"os"
)

// Config holds the provider configuration.
type Config struct {
	// VaultURL is the URL to the Azure Key Vault instance.
	// Format: https://<vault-name>.vault.azure.net/
	VaultURL string
}

// Load reads configuration from environment variables.
// It requires AZURE_KEYVAULT_URL to be set.
//
// Returns an error if AZURE_KEYVAULT_URL is not set.
func Load() (Config, error) {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")

	if vaultURL == "" {
		return Config{}, fmt.Errorf(
			"AZURE_KEYVAULT_URL is required",
		)
	}

	return Config{
		VaultURL: vaultURL,
	}, nil
}

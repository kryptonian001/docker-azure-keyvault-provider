package config

import (
	"fmt"
	"os"
)

type Config struct {
	VaultURL string
}

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

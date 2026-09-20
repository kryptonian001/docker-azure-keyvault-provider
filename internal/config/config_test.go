package config

import (
	"os"
	"testing"
)

// TestLoadWithValidConfig tests loading configuration with valid environment variables
func TestLoadWithValidConfig(t *testing.T) {
	// Set up environment variable
	vaultURL := "https://testvault.vault.azure.net/"
	os.Setenv("AZURE_KEYVAULT_URL", vaultURL)
	defer os.Unsetenv("AZURE_KEYVAULT_URL")

	config, err := Load()

	if err != nil {
		t.Errorf("Load() returned error %v, want nil", err)
	}

	if config.VaultURL != vaultURL {
		t.Errorf("Load() VaultURL = %v, want %v", config.VaultURL, vaultURL)
	}
}

// TestLoadWithMissingEnvVar tests loading configuration without required environment variable
func TestLoadWithMissingEnvVar(t *testing.T) {
	// Ensure the environment variable is not set
	os.Unsetenv("AZURE_KEYVAULT_URL")

	config, err := Load()

	if err == nil {
		t.Errorf("Load() returned nil error, want error about missing AZURE_KEYVAULT_URL")
	}

	if config.VaultURL != "" {
		t.Errorf("Load() VaultURL = %v, want empty string", config.VaultURL)
	}

	if !containsSubstring(err.Error(), "AZURE_KEYVAULT_URL is required") {
		t.Errorf("Load() error = %v, want to contain 'AZURE_KEYVAULT_URL is required'", err.Error())
	}
}

// TestLoadWithEmptyEnvVar tests loading configuration with empty environment variable
func TestLoadWithEmptyEnvVar(t *testing.T) {
	// Set environment variable to empty string
	os.Setenv("AZURE_KEYVAULT_URL", "")
	defer os.Unsetenv("AZURE_KEYVAULT_URL")

	config, err := Load()

	if err == nil {
		t.Errorf("Load() returned nil error, want error about empty AZURE_KEYVAULT_URL")
	}

	if config.VaultURL != "" {
		t.Errorf("Load() VaultURL = %v, want empty string", config.VaultURL)
	}
}

// TestLoadWithDifferentURLFormats tests loading configuration with various URL formats
func TestLoadWithDifferentURLFormats(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "standard vault URL",
			url:     "https://myvault.vault.azure.net/",
			wantErr: false,
		},
		{
			name:    "vault URL without trailing slash",
			url:     "https://myvault.vault.azure.net",
			wantErr: false,
		},
		{
			name:    "vault URL with different region",
			url:     "https://chinVault.vault.azure.cn/",
			wantErr: false,
		},
		{
			name:    "http URL (should still parse)",
			url:     "http://localhost:8200",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("AZURE_KEYVAULT_URL", tt.url)
			defer os.Unsetenv("AZURE_KEYVAULT_URL")

			config, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && config.VaultURL != tt.url {
				t.Errorf("Load() VaultURL = %v, want %v", config.VaultURL, tt.url)
			}
		})
	}
}

// TestLoadMultipleCalls tests that Load can be called multiple times
func TestLoadMultipleCalls(t *testing.T) {
	vaultURL1 := "https://vault1.vault.azure.net/"
	vaultURL2 := "https://vault2.vault.azure.net/"

	// First call
	os.Setenv("AZURE_KEYVAULT_URL", vaultURL1)
	config1, err1 := Load()
	if err1 != nil {
		t.Errorf("Load() first call returned error %v", err1)
	}

	// Change environment variable
	os.Setenv("AZURE_KEYVAULT_URL", vaultURL2)
	config2, err2 := Load()
	if err2 != nil {
		t.Errorf("Load() second call returned error %v", err2)
	}

	// Verify they're different
	if config1.VaultURL == config2.VaultURL {
		t.Errorf("Load() should return different values on different calls")
	}

	if config1.VaultURL != vaultURL1 {
		t.Errorf("Load() first call = %v, want %v", config1.VaultURL, vaultURL1)
	}

	if config2.VaultURL != vaultURL2 {
		t.Errorf("Load() second call = %v, want %v", config2.VaultURL, vaultURL2)
	}

	os.Unsetenv("AZURE_KEYVAULT_URL")
}

// containsSubstring is a helper function
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

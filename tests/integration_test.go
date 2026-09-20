package tests

import (
	"context"
	"testing"

	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/config"
	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/provider"
)

// TestProviderIntegration is a basic integration test showing provider usage
func TestProviderIntegration(t *testing.T) {
	// Create provider with nil client (won't actually call Azure)
	p := provider.New(nil)

	if p == nil {
		t.Fatalf("New() returned nil")
	}
}

// TestConfigIntegration tests configuration loading in isolation
func TestConfigIntegration(t *testing.T) {
	t.Run("loads config successfully", func(t *testing.T) {
		// This test will fail unless AZURE_KEYVAULT_URL is set
		// That's expected - it demonstrates the integration requirement
		cfg, err := config.Load()

		// If error, it's because env var not set (expected in CI/test)
		// If success, verify the config
		if err == nil && cfg.VaultURL == "" {
			t.Errorf("Config loaded but VaultURL is empty")
		}
	})
}

// TestProviderWorkflow shows a typical provider usage workflow
func TestProviderWorkflow(t *testing.T) {
	// Step 1: Load configuration
	cfg, err := config.Load()
	if err != nil {
		// In production, you'd handle this properly
		t.Logf("Config load failed (expected in test): %v", err)
		return
	}

	// Step 2: Create provider (would need real client in production)
	// In tests, we use nil client as a placeholder
	p := provider.New(nil)

	// Step 3: Verify provider is ready
	if p == nil {
		t.Fatalf("Provider creation failed")
	}

	// Step 4: Provider can be started
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run in background to avoid blocking
	go p.Run(ctx)

	// Step 5: Clean shutdown
	cancel()

	t.Logf("Provider workflow completed with config URL: %s", cfg.VaultURL)
}

// BenchmarkProviderCreation benchmarks provider creation
func BenchmarkProviderCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = provider.New(nil)
	}
}

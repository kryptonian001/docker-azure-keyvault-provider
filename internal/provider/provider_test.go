package provider

import (
	"context"
	"testing"
	"time"

	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/logging"
)

// TestParseSecretName tests the parseSecretName function with various inputs
func TestParseSecretName(t *testing.T) {
	tests := []struct {
		name      string
		requestID string
		want      string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid secret name",
			requestID: "azure/my-secret",
			want:      "my-secret",
			wantErr:   false,
		},
		{
			name:      "valid secret name with dashes",
			requestID: "azure/my-database-password",
			want:      "my-database-password",
			wantErr:   false,
		},
		{
			name:      "valid secret name with underscores",
			requestID: "azure/my_secret_123",
			want:      "my_secret_123",
			wantErr:   false,
		},
		{
			name:      "valid secret name with numbers",
			requestID: "azure/secret123",
			want:      "secret123",
			wantErr:   false,
		},
		{
			name:      "missing azure prefix",
			requestID: "my-secret",
			want:      "",
			wantErr:   true,
			errMsg:    "unsupported secret ID",
		},
		{
			name:      "wrong prefix",
			requestID: "aws/my-secret",
			want:      "",
			wantErr:   true,
			errMsg:    "unsupported secret ID",
		},
		{
			name:      "empty secret name",
			requestID: "azure/",
			want:      "",
			wantErr:   true,
			errMsg:    "Secret name cannot be empty",
		},
		{
			name:      "only prefix",
			requestID: "azure",
			want:      "",
			wantErr:   true,
			errMsg:    "unsupported secret ID",
		},
		{
			name:      "nested path not supported",
			requestID: "azure/folder/secret",
			want:      "",
			wantErr:   true,
			errMsg:    "nested paths are not supported",
		},
		{
			name:      "multiple nested paths",
			requestID: "azure/path/to/secret",
			want:      "",
			wantErr:   true,
			errMsg:    "nested paths are not supported",
		},
		{
			name:      "secret with spaces",
			requestID: "azure/my secret",
			want:      "my secret",
			wantErr:   false,
		},
		{
			name:      "secret with special chars",
			requestID: "azure/secret-with-special_chars.123",
			want:      "secret-with-special_chars.123",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSecretName(tt.requestID)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseSecretName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && !containsSubstring(err.Error(), tt.errMsg) {
				t.Errorf("parseSecretName() error message = %v, want to contain %v", err.Error(), tt.errMsg)
			}

			if got != tt.want {
				t.Errorf("parseSecretName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewProvider tests provider creation
func TestNewProvider(t *testing.T) {
	// Create a nil client to test
	provider := New(nil)

	if provider == nil {
		t.Errorf("New() returned nil, want provider instance")
	}

	if provider.client != nil {
		t.Errorf("New() with nil client should have nil client, got %v", provider.client)
	}
}

// TestNewProviderWithClient tests provider creation with a client
func TestNewProviderWithClient(t *testing.T) {
	// Note: In a real scenario, you'd use a mock client
	provider := New(nil)

	if provider == nil {
		t.Errorf("New() returned nil")
	}
}

// TestRunProvider tests the Run method behavior
func TestRunProvider(t *testing.T) {
	provider := New(nil)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Run the provider in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- provider.Run(ctx)
	}()

	// Give the goroutine time to start
	time.Sleep(10 * time.Millisecond)

	// Cancel the context
	cancel()

	// Wait for Run to return
	err := <-done

	if err != nil {
		t.Errorf("Run() returned error %v, want nil", err)
	}
}

// containsSubstring is a helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestNewProviderWithLogger tests provider creation with a logger
func TestNewProviderWithLogger(t *testing.T) {
	logger := &logging.MockLogger{}
	provider := NewWithLogger(nil, logger)

	if provider == nil {
		t.Errorf("NewWithLogger() returned nil, want provider instance")
	}

	if provider.client != nil {
		t.Errorf("NewWithLogger() with nil client should have nil client, got %v", provider.client)
	}

	if provider.log != logger {
		t.Errorf("NewWithLogger() logger not set correctly")
	}
}

// mockPattern is a test helper for the plugin.Pattern interface
type mockPattern struct {
	value string
}

func (m *mockPattern) String() string {
	return m.value
}

func (m *mockPattern) Dir() string {
	return ""
}

func (m *mockPattern) ExpandID(id interface{}) (interface{}, error) {
	return id, nil
}

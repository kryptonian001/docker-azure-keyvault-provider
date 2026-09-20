package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/docker/secrets-engine/plugin"
)

const secretPrefix = "azure/"

// AzureKeyVaultProvider implements the Docker Secrets Engine provider interface
// for Azure Key Vault. It retrieves secrets from Azure Key Vault based on secret IDs.
//
// Secret IDs must follow the format "azure/<secret-name>". The secret name cannot
// contain forward slashes or be empty.
type AzureKeyVaultProvider struct {
	client *azsecrets.Client
}

// New creates a new AzureKeyVaultProvider with the given Azure secrets client.
// The client is used to retrieve secrets from Azure Key Vault.
func New(client *azsecrets.Client) *AzureKeyVaultProvider {
	return &AzureKeyVaultProvider{
		client: client,
	}
}

// GetSecrets retrieves a secret from Azure Key Vault based on the provided pattern.
// The pattern must be a valid secret ID in the format "azure/<secret-name>".
//
// Returns an error if:
// - The secret ID is missing the "azure/" prefix
// - The secret name is empty
// - The secret name contains forward slashes (nested paths)
// - The secret does not exist in Azure Key Vault
// - Azure Key Vault returns an empty secret value
func (p *AzureKeyVaultProvider) GetSecrets(
	ctx context.Context,
	pattern plugin.Pattern,
) ([]plugin.Envelope, error) {

	requestID := pattern.String()

	secretName, err := parseSecretName(requestID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	response, err := p.client.GetSecret(
		ctx,
		secretName,
		"",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to retrieve secret %q from Azure Key Vault: %w",
			secretName,
			err,
		)
	}

	if response.Value == nil {
		return nil, fmt.Errorf(
			"Azure Key Vault returned an empty value for secret %q",
			secretName,
		)
	}

	id := plugin.MustParseID(requestID)

	return []plugin.Envelope{
		{
			ID:        id,
			Value:     []byte(*response.Value),
			CreatedAt: time.Now(),
		},
	}, nil
}

// Run starts the provider and keeps it running until the context is cancelled.
// This method blocks until ctx.Done() is signaled.
func (p *AzureKeyVaultProvider) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// parseSecretName extracts the secret name from a request ID.
// The request ID must be in the format "azure/<secret-name>".
//
// Returns an error if:
// - The request ID does not start with "azure/"
// - The secret name is empty
// - The secret name contains forward slashes (nested paths not supported)
func parseSecretName(requestID string) (string, error) {
	if !strings.HasPrefix(requestID, secretPrefix) {
		return "", fmt.Errorf(
			"unsupported secret ID %q: expected prefix %q",
			requestID,
			secretPrefix,
		)
	}

	secretName := strings.TrimPrefix(requestID, secretPrefix)

	if strings.Contains(secretName, "/") {
		return "", fmt.Errorf(
			"invalid secret name %q: nested paths are not supported",
			secretName,
		)
	}

	return secretName, nil
}

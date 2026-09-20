package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/docker/secrets-engine/plugin"
	"github.com/kryptonian001/docker-azure-keyvault-provider/internal/logging"
)

const secretPrefix = "azure/"

// AzureKeyVaultProvider implements the Docker Secrets Engine provider interface
// for Azure Key Vault. It retrieves secrets from Azure Key Vault based on secret IDs.
//
// Secret IDs must follow the format "azure/<secret-name>". The secret name cannot
// contain forward slashes or be empty.
type AzureKeyVaultProvider struct {
	client *azsecrets.Client
	log    logging.Logger
}

// New creates a new AzureKeyVaultProvider with the given Azure secrets client.
// The client is used to retrieve secrets from Azure Key Vault.
func New(client *azsecrets.Client) *AzureKeyVaultProvider {
	return &AzureKeyVaultProvider{
		client: client,
		log:    logging.NewNoOpLogger(),
	}
}

// NewWithLogger creates a new AzureKeyVaultProvider with the given Azure secrets client
// and logger for structured logging.
func NewWithLogger(client *azsecrets.Client, log logging.Logger) *AzureKeyVaultProvider {
	return &AzureKeyVaultProvider{
		client: client,
		log:    log,
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
	p.log.Debug("GetSecrets called", "requestID", requestID)

	secretName, err := parseSecretName(requestID)
	if err != nil {
		p.log.Error("failed to parse secret name", "error", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	p.log.Debug("retrieving secret from Azure Key Vault", "secretName", secretName)
	response, err := p.client.GetSecret(
		ctx,
		secretName,
		"",
		nil,
	)
	if err != nil {
		p.log.Error("failed to retrieve secret", "secretName", secretName, "error", err)
		return nil, fmt.Errorf(
			"failed to retrieve secret %q from Azure Key Vault: %w",
			secretName,
			err,
		)
	}

	if response.Value == nil {
		p.log.Error("Azure Key Vault returned empty value", "secretName", secretName)
		return nil, fmt.Errorf(
			"Azure Key Vault returned an empty value for secret %q",
			secretName,
		)
	}

	id := plugin.MustParseID(requestID)
	p.log.Info("successfully retrieved secret", "secretName", secretName)

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

	if secretName == "" {
		return "", fmt.Errorf("Secret name cannot be empty")
	}

	if strings.Contains(secretName, "/") {
		return "", fmt.Errorf(
			"invalid secret name %q: nested paths are not supported",
			secretName,
		)
	}

	return secretName, nil
}

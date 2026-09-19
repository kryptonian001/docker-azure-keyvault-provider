package provider

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/docker/secrets-engine/plugin"
)

type AzureKeyVaultProvider struct {
	client *azsecrets.Client
}

func New(client *azsecrets.Client) *AzureKeyVaultProvider {
	return &AzureKeyVaultProvider{
		client: client,
	}
}

func (p *AzureKeyVaultProvider) GetSecrets(
	ctx context.Context,
	pattern plugin.Pattern,
) ([]plugin.Envelope, error) {

	return nil, nil
}

func (p *AzureKeyVaultProvider) Run(ctx context.Context) error {
	<-ctx.Done()

	return nil
}

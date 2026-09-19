package plugin

import (
	"context"

	"github.com/docker/secrets-engine/plugin"
)

type Plugin interface {
	GetSecrets(
		context.Context,
		plugin.Pattern,
	) ([]plugin.Envelope, error)

	Run(context.Context) error
}

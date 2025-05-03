package app

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
)

type Registry struct {
}

func NewRepo(
	ctx context.Context,
	version string,
	creds *configuration.Credentials,
) (*Registry, error) {

	return InitRepoRegistry()
}

func InitRepoRegistry() (*Registry, error) {
	return &Registry{}, nil
}

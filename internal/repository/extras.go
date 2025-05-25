package repository

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
)

type IExtras interface {
	GetAll(ctx context.Context) ([]entities.Extra, error)
	Add(ctx context.Context, extras []entities.Extra) error
	Update(ctx context.Context, extra entities.Extra) error
	Delete(ctx context.Context, id int) error
}

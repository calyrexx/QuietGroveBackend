package repository

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
)

type IHouses interface {
	GetAll(ctx context.Context) ([]entities.House, error)
	GetOne(ctx context.Context, id int) (entities.House, error)
	Add(ctx context.Context, house entities.House) error
	Update(ctx context.Context, house entities.House) error
	Delete(ctx context.Context, id int) error
}

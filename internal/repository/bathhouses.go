package repository

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
)

type IBathhouses interface {
	GetAll(ctx context.Context) ([]entities.Bathhouse, error)
	GetByHouse(ctx context.Context, houseID int) ([]entities.Bathhouse, error)
	Add(ctx context.Context, bathhouses []entities.Bathhouse) error
	Update(ctx context.Context, bathhouse entities.Bathhouse) error
	Delete(ctx context.Context, id int) error
}

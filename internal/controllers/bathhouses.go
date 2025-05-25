package controllers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
)

type IBathhousesUseCase interface {
	GetAll(ctx context.Context) ([]entities.Bathhouse, error)
	GetByHouse(ctx context.Context, houseID int) ([]entities.Bathhouse, error)
	Add(ctx context.Context, bhs []entities.Bathhouse) error
	Update(ctx context.Context, bh entities.Bathhouse) error
	Delete(ctx context.Context, id int) error
}

type BathhousesDependencies struct {
	UseCase IBathhousesUseCase
}

type Bathhouses struct {
	useCase IBathhousesUseCase
}

func NewBathhouses(d *BathhousesDependencies) (*Bathhouses, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Bathhouses Controller", "usecase", "nil")
	}
	return &Bathhouses{
		useCase: d.UseCase,
	}, nil
}

func (c *Bathhouses) GetAll(ctx context.Context) ([]entities.Bathhouse, error) {
	return c.useCase.GetAll(ctx)
}

func (c *Bathhouses) GetByHouse(ctx context.Context, houseID int) ([]entities.Bathhouse, error) {
	return c.useCase.GetByHouse(ctx, houseID)
}

func (c *Bathhouses) Add(ctx context.Context, bhs []entities.Bathhouse) error {
	return c.useCase.Add(ctx, bhs)
}

func (c *Bathhouses) Update(ctx context.Context, bh entities.Bathhouse) error {
	return c.useCase.Update(ctx, bh)
}

func (c *Bathhouses) Delete(ctx context.Context, id int) error {
	return c.useCase.Delete(ctx, id)
}

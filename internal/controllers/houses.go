package controllers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
)

type IHousesUseCase interface {
	GetAll(ctx context.Context) ([]entities.House, error)
	Add(ctx context.Context, house entities.House) error
	Update(ctx context.Context, house entities.House) error
	Delete(ctx context.Context, houseID int) error
}

type HousesDependencies struct {
	UseCase IHousesUseCase
}

type Houses struct {
	useCase IHousesUseCase
}

func NewHouses(d *HousesDependencies) (*Houses, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Houses UseCase", "whole", "nil")
	}
	return &Houses{
		useCase: d.UseCase,
	}, nil
}

func (c *Houses) GetAll(ctx context.Context) ([]entities.House, error) {
	return c.useCase.GetAll(ctx)
}

func (c *Houses) Add(ctx context.Context, house entities.House) error {
	return c.useCase.Add(ctx, house)
}

func (c *Houses) Update(ctx context.Context, house entities.House) error {
	return c.useCase.Update(ctx, house)
}

func (c *Houses) Delete(ctx context.Context, houseID int) error {
	return c.useCase.Delete(ctx, houseID)
}

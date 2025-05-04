package controllers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api/handlers"
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

func (c *Houses) GetAll(ctx context.Context) ([]handlers.House, error) {
	res, err := c.useCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return c.convertEntitiesToHouses(res), nil
}

func (c *Houses) convertEntitiesToHouses(enities []entities.House) []handlers.House {
	res := make([]handlers.House, 0, len(enities))
	for _, entity := range enities {
		res = append(res, c.convertEntityToHouse(entity))
	}
	return res
}

func (c *Houses) convertEntityToHouse(entity entities.House) handlers.House {
	return handlers.House{
		ID:            entity.ID,
		Name:          entity.Name,
		Description:   entity.Description,
		Capacity:      entity.Capacity,
		BasePrice:     entity.BasePrice,
		Images:        entity.Images,
		CheckInFrom:   entity.CheckInFrom,
		CheckOutUntil: entity.CheckOutUntil,
	}
}

func (c *Houses) convertHouseToEntity(house handlers.House) entities.House {
	return entities.House{
		ID:            house.ID,
		Name:          house.Name,
		Description:   house.Description,
		Capacity:      house.Capacity,
		BasePrice:     house.BasePrice,
		Images:        house.Images,
		CheckInFrom:   house.CheckInFrom,
		CheckOutUntil: house.CheckOutUntil,
	}
}

func (c *Houses) Add(ctx context.Context, house handlers.House) error {
	return c.useCase.Add(ctx, c.convertHouseToEntity(house))
}

func (c *Houses) Update(ctx context.Context, house entities.House) error {
	return c.useCase.Update(ctx, house)
}

func (c *Houses) Delete(ctx context.Context, houseID int) error {
	return c.useCase.Delete(ctx, houseID)
}

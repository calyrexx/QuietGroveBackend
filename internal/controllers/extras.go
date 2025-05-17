package controllers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api/handlers"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
)

type IExtrasUseCase interface {
	GetAll(ctx context.Context) ([]entities.Extra, error)
	Add(ctx context.Context, extra []entities.Extra) error
	Update(ctx context.Context, extra entities.Extra) error
	Delete(ctx context.Context, extraID int) error
}

type ExtrasDependencies struct {
	UseCase IExtrasUseCase
}

type Extras struct {
	useCase IExtrasUseCase
}

func NewExtras(d *ExtrasDependencies) (*Extras, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Extras UseCase", "whole", "nil")
	}
	return &Extras{
		useCase: d.UseCase,
	}, nil
}

func (c *Extras) GetAll(ctx context.Context) ([]handlers.Extra, error) {
	res, err := c.useCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return c.convertEntitiesToExtras(res), nil
}

func (c *Extras) Add(ctx context.Context, extras []handlers.Extra) error {
	if len(extras) == 0 {
		return nil
	}

	return c.useCase.Add(ctx, c.convertExtrasToEntity(extras))
}

func (c *Extras) Update(ctx context.Context, extra handlers.Extra) error {
	return c.useCase.Update(ctx, c.convertExtraToEntity(extra))
}

func (c *Extras) Delete(ctx context.Context, extraID int) error {
	return c.useCase.Delete(ctx, extraID)
}

func (c *Extras) convertEntitiesToExtras(entities []entities.Extra) []handlers.Extra {
	res := make([]handlers.Extra, 0, len(entities))
	for _, entity := range entities {
		res = append(res, c.convertEntityToExtra(entity))
	}
	return res
}

func (c *Extras) convertEntityToExtra(entity entities.Extra) handlers.Extra {
	return handlers.Extra{
		ID:          entity.ID,
		Name:        entity.Name,
		Text:        entity.Text,
		Description: entity.Description,
		BasePrice:   entity.BasePrice,
		Images:      entity.Images,
	}
}

func (c *Extras) convertExtrasToEntity(extras []handlers.Extra) []entities.Extra {
	resp := make([]entities.Extra, 0, len(extras))
	for _, extra := range extras {
		resp = append(resp, c.convertExtraToEntity(extra))
	}
	return resp
}

func (c *Extras) convertExtraToEntity(extra handlers.Extra) entities.Extra {
	return entities.Extra{
		ID:          extra.ID,
		Name:        extra.Name,
		Text:        extra.Text,
		Description: extra.Description,
		BasePrice:   extra.BasePrice,
		Images:      extra.Images,
	}
}

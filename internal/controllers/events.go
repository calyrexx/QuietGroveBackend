package controllers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/api/handlers"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
)

type IEventsUseCase interface {
	NewApplication(ctx context.Context, req entities.NewApplication) error
}

type EventsDependencies struct {
	UseCase IEventsUseCase
}

type Events struct {
	useCase IEventsUseCase
}

func NewEvents(d *EventsDependencies) (*Events, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Events UseCase", "whole", "nil")
	}

	return &Events{
		useCase: d.UseCase,
	}, nil
}

func (c *Events) NewApplication(ctx context.Context, req handlers.EventsNewApplication) error {
	request := c.convertNewApplication(req)

	return c.useCase.NewApplication(ctx, request)
}

func (c *Events) convertNewApplication(req handlers.EventsNewApplication) entities.NewApplication {
	return entities.NewApplication{
		Name:        req.Name,
		Phone:       req.Phone,
		CheckIn:     req.CheckIn,
		GuestsCount: req.GuestsCount,
	}
}

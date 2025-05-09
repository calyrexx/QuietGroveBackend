package controllers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/usecases"
)

type IReservationsUseCase interface {
	CreateReservation(ctx context.Context, req usecases.CreateReservationRequest) (entities.Reservation, error)
}

type ReservationsDependencies struct {
	UseCase IReservationsUseCase
}

type Reservations struct {
	useCase IReservationsUseCase
}

func NewReservations(d *ReservationsDependencies) (*Reservations, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Reservations UseCase", "whole", "nil")
	}
	return &Reservations{
		useCase: d.UseCase,
	}, nil
}

func (c *Reservations) BookAHouse(ctx context.Context, req usecases.CreateReservationRequest) (entities.Reservation, error) {
	response, err := c.useCase.CreateReservation(ctx, req)
	if err != nil {
		return entities.Reservation{}, err
	}
	return response, nil
}

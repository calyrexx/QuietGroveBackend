package controllers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
)

type IReservationsUseCase interface {
	BookAHouse(ctx context.Context, request string) error
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

func (c *Reservations) BookAHouse(ctx context.Context, request string) (string, error) {
	if err := c.useCase.BookAHouse(ctx, request); err != nil {
		return "", err
	}
	return "ok", nil
}

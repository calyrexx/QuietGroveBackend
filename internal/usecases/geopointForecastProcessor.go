package usecases

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/sirupsen/logrus"
)

type (
	IReservationsRepository interface {
	}

	ReservationsDependencies struct {
		Repo   IReservationsRepository
		Logger logrus.FieldLogger
	}
	Reservations struct {
		repo   IReservationsRepository
		logger logrus.FieldLogger
	}
)

func NewReservations(d *ReservationsDependencies) (*Reservations, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservations", "whole", "nil")
	}

	logger := d.Logger.WithField("Usecases", "Reservations")

	return &Reservations{
		repo:   d.Repo,
		logger: logger,
	}, nil
}

func (u *Reservations) BookAHouse(ctx context.Context, req string) error {
	return nil
}

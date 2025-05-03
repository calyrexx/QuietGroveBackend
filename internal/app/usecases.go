package app

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/usecases"
	"github.com/sirupsen/logrus"
)

type Usecases struct {
	reservations *usecases.Reservations
	houses       *usecases.Houses
}

func NewUsecases(
	logger logrus.FieldLogger,
	repo *Registry,
) (*Usecases, error) {

	reservationsUsecase, err := usecases.NewReservations(&usecases.ReservationsDependencies{
		Repo:   repo,
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	housesUsecase, err := usecases.NewHouses(&usecases.HousesDependencies{
		Repo:   repo.Houses,
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	return &Usecases{
		reservations: reservationsUsecase,
		houses:       housesUsecase,
	}, nil
}

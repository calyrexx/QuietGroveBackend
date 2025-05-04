package app

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/controllers"
	"github.com/sirupsen/logrus"
)

type Controllers struct {
	Reservations *controllers.Reservations
	Houses       *controllers.Houses
	Extras       *controllers.Extras
}

func NewControllers(
	logger logrus.FieldLogger,
	usecases *Usecases,
) (*Controllers, error) {

	reservationsController, err := controllers.NewReservations(&controllers.ReservationsDependencies{
		UseCase: usecases.reservations,
	})
	if err != nil {
		return nil, err
	}

	housesController, err := controllers.NewHouses(&controllers.HousesDependencies{
		UseCase: usecases.houses,
	})
	if err != nil {
		return nil, err
	}

	extrasController, err := controllers.NewExtras(&controllers.ExtrasDependencies{
		UseCase: usecases.extras,
	})
	if err != nil {
		return nil, err
	}

	return &Controllers{
		Reservations: reservationsController,
		Houses:       housesController,
		Extras:       extrasController,
	}, nil
}

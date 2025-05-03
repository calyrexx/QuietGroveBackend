package app

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/controllers"
	"github.com/sirupsen/logrus"
)

type Controllers struct {
	Reservations *controllers.Reservations
}

func NewControllers(
	logger logrus.FieldLogger,
	usecases *Usecases,
) (*Controllers, error) {

	reservationsController, err := controllers.NewReservations(&controllers.ReservationsDependencies{
		UseCase: usecases.reservations,
	})
	if err != nil {
		logger.Fatalf("controllers.NewReservations init error: %v", err)
	}

	return &Controllers{
		Reservations: reservationsController,
	}, nil
}

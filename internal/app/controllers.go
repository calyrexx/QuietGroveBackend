package app

import (
	"github.com/calyrexx/QuietGrooveBackend/internal/controllers"
	"log/slog"
)

type Controllers struct {
	Reservations *controllers.Reservations
	Houses       *controllers.Houses
	Bathhouses   *controllers.Bathhouses
	Extras       *controllers.Extras
	Verification *controllers.Verification
	Events       *controllers.Events
}

func NewControllers(
	logger *slog.Logger,
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

	bathhousesController, err := controllers.NewBathhouses(&controllers.BathhousesDependencies{
		UseCase: usecases.bathhouses,
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

	verificationController, err := controllers.NewVerification(&controllers.VerificationDependencies{
		UseCase: usecases.verification,
	})
	if err != nil {
		return nil, err
	}

	eventsController, err := controllers.NewEvents(&controllers.EventsDependencies{
		UseCase: usecases.events,
	})
	if err != nil {
		return nil, err
	}

	return &Controllers{
		Reservations: reservationsController,
		Houses:       housesController,
		Bathhouses:   bathhousesController,
		Extras:       extrasController,
		Verification: verificationController,
		Events:       eventsController,
	}, nil
}

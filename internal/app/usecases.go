package app

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/Calyr3x/QuietGrooveBackend/internal/integrations/telegram"
	"github.com/Calyr3x/QuietGrooveBackend/internal/usecases"
	"github.com/sirupsen/logrus"
)

type Usecases struct {
	reservations *usecases.Reservation
	houses       *usecases.Houses
	extras       *usecases.Extras
}

func NewUsecases(
	logger logrus.FieldLogger,
	config *configuration.Config,
	repo *Registry,
	tgBot *telegram.TGNotifier,
) (*Usecases, error) {

	reservationsUsecase, err := usecases.NewReservation(&usecases.ReservationDependencies{
		ReservationRepo: repo.Reservations,
		GuestRepo:       repo.Guests,
		HouseRepo:       repo.Houses,
		PCoefs:          config.PriceCoefficients,
		Logger:          logger,
		Notifier:        tgBot,
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

	extrasUsecase, err := usecases.NewExtras(&usecases.ExtrasDependencies{
		Repo:   repo.Extras,
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	return &Usecases{
		reservations: reservationsUsecase,
		houses:       housesUsecase,
		extras:       extrasUsecase,
	}, nil
}

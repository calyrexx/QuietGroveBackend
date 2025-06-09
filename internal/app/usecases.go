package app

import (
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/QuietGrooveBackend/internal/integrations/telegram"
	"github.com/calyrexx/QuietGrooveBackend/internal/usecases"
	"log/slog"
	"time"
)

type Usecases struct {
	reservations *usecases.Reservation
	houses       *usecases.Houses
	bathhouses   *usecases.Bathhouses
	extras       *usecases.Extras
	verification *usecases.Verification
	events       *usecases.Events
}

func NewUsecases(
	logger *slog.Logger,
	config *configuration.Config,
	repo *Registry,
	tgBot *telegram.Adapter,
) (*Usecases, error) {

	reservationsUsecase, err := usecases.NewReservation(&usecases.ReservationDependencies{
		ReservationRepo: repo.Reservations,
		GuestRepo:       repo.Guests,
		HouseRepo:       repo.Houses,
		BathhouseRepo:   repo.Bathhouses,
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

	bathhousesUsecase, err := usecases.NewBathhouses(&usecases.BathhousesDependencies{
		Repo:   repo.Bathhouses,
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

	verificationUsecase, err := usecases.NewVerification(&usecases.VerificationDependencies{
		Repo:       repo.Verification,
		GuestsRepo: repo.Guests,
		TTL:        time.Hour,
	})
	if err != nil {
		return nil, err
	}

	eventsUsecase, err := usecases.NewEvents(&usecases.EventsDependencies{
		Logger:   logger,
		Notifier: tgBot,
	})
	if err != nil {
		return nil, err
	}

	return &Usecases{
		reservations: reservationsUsecase,
		houses:       housesUsecase,
		bathhouses:   bathhousesUsecase,
		extras:       extrasUsecase,
		verification: verificationUsecase,
		events:       eventsUsecase,
	}, nil
}

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
	extras       *usecases.Extras
	verification *usecases.Verification
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

	verificationUsecase, err := usecases.NewVerification(&usecases.VerificationDependencies{
		Repo: repo.Verification,
		TTL:  time.Hour,
	})
	if err != nil {
		return nil, err
	}

	return &Usecases{
		reservations: reservationsUsecase,
		houses:       housesUsecase,
		extras:       extrasUsecase,
		verification: verificationUsecase,
	}, nil
}

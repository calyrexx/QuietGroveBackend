package app

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Registry struct {
	Reservations repository.IReservations
	Houses       repository.IHouses
	Extras       repository.IExtras
	Guests       repository.IGuests
	Verification repository.IVerification
}

func NewRepo(ctx context.Context, creds *configuration.Credentials) (*Registry, error) {
	postgresConnect, err := postgres.NewPostgres(ctx, creds.Postgres)
	if err != nil {
		return nil, err
	}

	return InitRepoRegistry(postgresConnect)
}

func InitRepoRegistry(postgresConnect *pgxpool.Pool) (*Registry, error) {
	reservationsRepo := postgres.NewReservationsRepo(postgresConnect)
	housesRepo := postgres.NewHousesRepo(postgresConnect)
	extrasRepo := postgres.NewExtrasRepo(postgresConnect)
	guestsRepo := postgres.NewGuestsRepo(postgresConnect)
	verificationRepo := postgres.NewVerificationRepo(postgresConnect)

	return &Registry{
		Reservations: reservationsRepo,
		Houses:       housesRepo,
		Extras:       extrasRepo,
		Guests:       guestsRepo,
		Verification: verificationRepo,
	}, nil
}

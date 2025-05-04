package app

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Registry struct {
	Reservations repository.IReservations
	Houses       repository.IHouses
	Extras       repository.IExtras
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

	return &Registry{
		Reservations: reservationsRepo,
		Houses:       housesRepo,
		Extras:       extrasRepo,
	}, nil
}

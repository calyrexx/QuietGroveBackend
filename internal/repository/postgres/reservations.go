package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationsRepo struct {
	pool *pgxpool.Pool
}

func NewReservationsRepo(pool *pgxpool.Pool) *ReservationsRepo {
	return &ReservationsRepo{pool: pool}
}

func (r *ReservationsRepo) BookAHouse(ctx context.Context, req string) error {
	return nil
}

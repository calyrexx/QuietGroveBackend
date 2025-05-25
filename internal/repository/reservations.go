package repository

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
)

type IReservations interface {
	GetAvailableHouses(ctx context.Context, req entities.GetAvailableHouses) ([]int, error)
	CheckAvailability(ctx context.Context, req entities.CheckAvailability) (bool, error)
	GetPrice(ctx context.Context, houseID int, extras []entities.ReservationExtra, bathhouse []entities.BathhouseReservation) (entities.GetPrice, error)
	Create(ctx context.Context, reservation entities.Reservation) error
}

package repository

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
)

type IReservations interface {
	GetAvailableHouses(ctx context.Context, req entities.GetAvailableHouses) ([]int, error)
	CheckAvailability(ctx context.Context, req entities.CheckAvailability) (bool, error)
	GetPrice(ctx context.Context, houseID int, extras []entities.ReservationExtra) (entities.GetPrice, error)
	Create(ctx context.Context, reservation entities.Reservation) error
}

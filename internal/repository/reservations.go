package repository

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
)

const DateFormat = "2006-01-02"

type IReservations interface {
	CheckAvailability(ctx context.Context, req entities.CheckAvailability) (bool, error)
	GetPrice(ctx context.Context, houseID int, extras []entities.ReservationExtra) (entities.GetPrice, error)
	Create(ctx context.Context, reservation entities.Reservation) error
}

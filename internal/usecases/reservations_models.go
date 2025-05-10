package usecases

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"time"
)

type CreateReservationRequest struct {
	HouseID     int
	Guest       entities.Guest
	CheckIn     time.Time
	CheckOut    time.Time
	GuestsCount int
	Extras      []entities.ReservationExtra
}

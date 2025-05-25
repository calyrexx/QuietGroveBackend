package usecases

import (
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
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

type GetAvailableHousesResponse struct {
	ID            int
	Name          string
	Description   string
	Capacity      int
	BasePrice     int
	TotalPrice    int
	Images        []string
	CheckInFrom   string
	CheckOutUntil string
}

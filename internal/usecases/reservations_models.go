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
	Bathhouse   []entities.BathhouseReservation
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
	Bathhouses    []BathhouseSlots
}

type BathhouseSlots struct {
	TypeID     int
	Name       string
	Slots      []BathhouseDateSlots
	FillOption []BathhouseFillOption
}

type BathhouseDateSlots struct {
	Date string
	Time []BathhouseTimeSlots
}

type BathhouseTimeSlots struct {
	TimeFrom string
	TimeTo   string
}

type BathhouseFillOption struct {
	ID          int
	Name        string
	Description string
	Price       int
}

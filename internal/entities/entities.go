package entities

import (
	"github.com/google/uuid"
	"time"
)

const DateFormat = "2006-01-02"

type (
	House struct {
		ID            int
		Name          string
		Description   string
		Capacity      int
		BasePrice     int
		Images        []string
		CheckInFrom   string
		CheckOutUntil string
	}

	Guest struct {
		Name  string
		Email string
		Phone string
	}

	Reservation struct {
		HouseID     int
		GuestUUID   uuid.UUID
		CheckIn     time.Time // [checkIn, checkOut)
		CheckOut    time.Time
		GuestsCount int
		Status      string
		TotalPrice  int
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Extras      []ReservationExtra
	}

	ReservationExtra struct {
		ExtraID  int
		Quantity int
		Amount   int
	}

	Extra struct {
		ID          int
		Name        string
		Text        string
		Description string
		BasePrice   int
		Images      []string
	}

	Payment struct {
		UUID            uuid.UUID
		ReservationUUID uuid.UUID
		Amount          int
		Currency        string
		Method          string
		Status          string // enum PaymentStatus
		GatewayTxID     string
		PaidAt          time.Time
	}

	GetPrice struct {
		House  int
		Extras int
	}

	GetAvailableHouses struct {
		CheckIn     time.Time
		CheckOut    time.Time
		GuestsCount int
	}

	CheckAvailability struct {
		HouseId  int
		CheckIn  time.Time
		CheckOut time.Time
	}
)

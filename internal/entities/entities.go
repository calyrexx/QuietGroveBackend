package entities

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type House struct {
	ID          uint
	Name        string
	Description string
	Capacity    uint
	BasePrice   int
	Images      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Guest struct {
	UUID      uuid.UUID
	Name      string
	Email     string
	Phone     string
	CreatedAt time.Time
}

type Reservation struct {
	UUID        uuid.UUID
	HouseID     uint
	GuestUUID   uuid.UUID
	Stay        pgtype.Range[pgtype.Date] // [checkIn, checkOut)
	GuestsCount uint
	Status      string // enum ReservationStatus
	TotalPrice  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Extras      []ReservationExtra
}

type ReservationExtra struct {
	ExtraID  int
	Quantity uint
	Amount   uint
}

type Payment struct {
	UUID            uuid.UUID
	ReservationUUID uuid.UUID
	Amount          int
	Currency        string
	Method          string
	Status          string // enum PaymentStatus
	GatewayTxID     string
	PaidAt          time.Time
}

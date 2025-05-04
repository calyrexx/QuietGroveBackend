package entities

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type House struct {
	ID            int
	Name          string
	Description   string
	Capacity      int
	BasePrice     int
	Images        []string
	CheckInFrom   string
	CheckOutUntil string
}

type Guest struct {
	UUID  uuid.UUID
	Name  string
	Email string
	Phone string
}

type Reservation struct {
	UUID        uuid.UUID
	HouseID     int
	GuestUUID   uuid.UUID
	Stay        pgtype.Range[pgtype.Date] // [checkIn, checkOut)
	GuestsCount int
	Status      string // enum ReservationStatus
	TotalPrice  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Extras      []ReservationExtra
}

type ReservationExtra struct {
	ExtraID  int
	Quantity int
	Amount   int
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

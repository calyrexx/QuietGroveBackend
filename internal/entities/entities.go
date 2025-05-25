package entities

import (
	"github.com/google/uuid"
	"time"
)

const (
	VerifPending  VerificationStatus = "pending"
	VerifApproved VerificationStatus = "approved"
	VerifExpired  VerificationStatus = "expired"
)

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
		Bathhouse   []BathhouseReservation
	}

	ReservationCreatedMessage struct {
		House       string
		GuestName   string
		GuestPhone  string
		CheckIn     time.Time // [checkIn, checkOut)
		CheckOut    time.Time
		GuestsCount int
		TotalPrice  int
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
		House     int
		Extras    int
		Bathhouse int
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

	VerificationStatus string

	Verification struct {
		ID         int64
		Code       string
		Email      string
		Phone      string
		Name       string
		TgUserID   *int64
		Status     VerificationStatus
		CreatedAt  time.Time
		VerifiedAt *time.Time
		ExpiresAt  time.Time
	}

	BathhouseReservation struct {
		TypeID       int
		Date         string
		TimeFrom     string
		TimeTo       string
		FillOptionID int
	}

	Bathhouse struct {
		ID          int
		HouseID     int
		Name        string
		Price       int
		Description string
		Images      []string
		FillOptions []BathhouseFillOption
	}

	BathhouseFillOption struct {
		ID          int
		BathhouseID int
		Name        string
		Image       string
		Description string
		Price       int
	}
)

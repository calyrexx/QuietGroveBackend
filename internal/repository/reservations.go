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
	GetDetailsByUUID(ctx context.Context, telegramID int64, uuid string) (entities.ReservationMessage, error)
	GetByTelegramID(ctx context.Context, telegramID int64) ([]entities.ReservationMessage, error)
	Cancel(ctx context.Context, userTgId int64, reservationUUID string) error
	GetAllConfirmed(ctx context.Context) ([]entities.ReservationUpdateStatus, error)
	UpdateStatuses(ctx context.Context, reservations []entities.ReservationUpdateStatus) error
	GetAllForReminder(ctx context.Context) ([]entities.ReservationReminderNotification, error)
}

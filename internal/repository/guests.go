package repository

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/google/uuid"
)

type IGuests interface {
	Get(ctx context.Context, guest entities.Guest) (Guest, error)
	Create(ctx context.Context, guest entities.Guest) error
}

type Guest struct {
	UUID  uuid.UUID
	Name  string
	Phone string
	Email string
	TgId  int64
}

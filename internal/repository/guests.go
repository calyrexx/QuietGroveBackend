package repository

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/google/uuid"
)

type IGuests interface {
	FindOrCreate(ctx context.Context, guest entities.Guest) (Guest, error)
}

type Guest struct {
	UUID  uuid.UUID
	Name  string
	Phone string
	Email string
}

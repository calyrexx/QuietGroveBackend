package repository

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
)

type IVerification interface {
	Create(ctx context.Context, v entities.Verification) error
	GetByCode(ctx context.Context, code string) (entities.Verification, error)
	Approve(ctx context.Context, id string, tgUserID int64) error
}

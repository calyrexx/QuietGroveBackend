package repository

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
)

type IVerification interface {
	Create(ctx context.Context, v entities.Verification) (int64, error)
	GetByCode(ctx context.Context, code string) (entities.Verification, error)
	Approve(ctx context.Context, id int64, tgUserID int64) error
}

package usecases

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"github.com/calyrexx/zeroslog"
	"log/slog"
)

type BathhousesDependencies struct {
	Repo   repository.IBathhouses
	Logger *slog.Logger
}
type Bathhouses struct {
	repo   repository.IBathhouses
	logger *slog.Logger
}

func NewBathhouses(d *BathhousesDependencies) (*Bathhouses, error) {
	logger := d.Logger.With(zeroslog.UsecaseKey, "Bathhouses")
	return &Bathhouses{
		repo:   d.Repo,
		logger: logger,
	}, nil
}

func (u *Bathhouses) GetAll(ctx context.Context) ([]entities.Bathhouse, error) {
	return u.repo.GetAll(ctx)
}

func (u *Bathhouses) Add(ctx context.Context, bhs []entities.Bathhouse) error {
	return u.repo.Add(ctx, bhs)
}

func (u *Bathhouses) Update(ctx context.Context, bh entities.Bathhouse) error {
	return u.repo.Update(ctx, bh)
}

func (u *Bathhouses) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

package usecases

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"github.com/calyrexx/zeroslog"
	"log/slog"
)

type (
	ExtrasDependencies struct {
		Repo   repository.IExtras
		Logger *slog.Logger
	}
	Extras struct {
		repo   repository.IExtras
		logger *slog.Logger
	}
)

func NewExtras(d *ExtrasDependencies) (*Extras, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Extras", "whole", "nil")
	}

	logger := d.Logger.With(zeroslog.UsecaseKey, "Extras")

	return &Extras{
		repo:   d.Repo,
		logger: logger,
	}, nil
}

func (u *Extras) GetAll(ctx context.Context) ([]entities.Extra, error) {
	return u.repo.GetAll(ctx)
}

func (u *Extras) Add(ctx context.Context, extras []entities.Extra) error {
	return u.repo.Add(ctx, extras)
}

func (u *Extras) Update(ctx context.Context, extra entities.Extra) error {
	return u.repo.Update(ctx, extra)
}

func (u *Extras) Delete(ctx context.Context, extraID int) error {
	return u.repo.Delete(ctx, extraID)
}

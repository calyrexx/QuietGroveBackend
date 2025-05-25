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
	HousesDependencies struct {
		Repo   repository.IHouses
		Logger *slog.Logger
	}
	Houses struct {
		repo   repository.IHouses
		logger *slog.Logger
	}
)

func NewHouses(d *HousesDependencies) (*Houses, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Houses", "whole", "nil")
	}

	logger := d.Logger.With(zeroslog.UsecaseKey, "Houses")

	return &Houses{
		repo:   d.Repo,
		logger: logger,
	}, nil
}

func (u *Houses) GetAll(ctx context.Context) ([]entities.House, error) {
	return u.repo.GetAll(ctx)
}

func (u *Houses) Add(ctx context.Context, houses []entities.House) error {
	return u.repo.Add(ctx, houses)
}

func (u *Houses) Update(ctx context.Context, house entities.House) error {
	return u.repo.Update(ctx, house)
}

func (u *Houses) Delete(ctx context.Context, houseID int) error {
	return u.repo.Delete(ctx, houseID)
}

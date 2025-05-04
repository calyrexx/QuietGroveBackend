package usecases

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository"
	"github.com/sirupsen/logrus"
)

type (
	HousesDependencies struct {
		Repo   repository.IHouses
		Logger logrus.FieldLogger
	}
	Houses struct {
		repo   repository.IHouses
		logger logrus.FieldLogger
	}
)

func NewHouses(d *HousesDependencies) (*Houses, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Houses", "whole", "nil")
	}

	logger := d.Logger.WithField("Usecases", "Houses")

	return &Houses{
		repo:   d.Repo,
		logger: logger,
	}, nil
}

func (u *Houses) GetAll(ctx context.Context) ([]entities.House, error) {
	return u.repo.GetAll(ctx)
}

func (u *Houses) Add(ctx context.Context, house entities.House) error {
	return u.repo.Add(ctx, house)
}

func (u *Houses) Update(ctx context.Context, house entities.House) error {
	return u.repo.Update(ctx, house)
}

func (u *Houses) Delete(ctx context.Context, houseID int) error {
	return u.repo.Delete(ctx, houseID)
}

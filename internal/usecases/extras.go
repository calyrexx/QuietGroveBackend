package usecases

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository"
	"github.com/sirupsen/logrus"
)

type (
	ExtrasDependencies struct {
		Repo   repository.IExtras
		Logger logrus.FieldLogger
	}
	Extras struct {
		repo   repository.IExtras
		logger logrus.FieldLogger
	}
)

func NewExtras(d *ExtrasDependencies) (*Extras, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Extras", "whole", "nil")
	}

	logger := d.Logger.WithField("Usecases", "Extras")

	return &Extras{
		repo:   d.Repo,
		logger: logger,
	}, nil
}

func (u *Extras) GetAll(ctx context.Context) ([]entities.Extra, error) {
	return u.repo.GetAll(ctx)
}

func (u *Extras) Add(ctx context.Context, extra entities.Extra) error {
	return u.repo.Add(ctx, extra)
}

func (u *Extras) Update(ctx context.Context, extra entities.Extra) error {
	return u.repo.Update(ctx, extra)
}

func (u *Extras) Delete(ctx context.Context, extraID int) error {
	return u.repo.Delete(ctx, extraID)
}

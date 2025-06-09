package usecases

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/zeroslog"
	"log/slog"
)

type (
	EventsNotifier interface {
		NewApplicationForEvent(res entities.NewApplication) error
	}

	EventsDependencies struct {
		Logger   *slog.Logger
		Notifier EventsNotifier
	}

	Events struct {
		logger   *slog.Logger
		notifier EventsNotifier
	}
)

func NewEvents(d *EventsDependencies) (*Events, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Events", "whole", "nil")
	}
	if d.Notifier == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Events", "Notifier", "nil")
	}

	logger := d.Logger.With(zeroslog.UsecaseKey, "Events")

	return &Events{
		logger:   logger,
		notifier: d.Notifier,
	}, nil
}

func (e *Events) NewApplication(ctx context.Context, req entities.NewApplication) error {
	if err := e.notifier.NewApplicationForEvent(req); err != nil {
		return err
	}

	return nil
}

package handlers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/api"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"log/slog"
	"net/http"
)

type IEventsController interface {
	NewApplication(ctx context.Context, req EventsNewApplication) error
}

type EventsDependencies struct {
	Controller IEventsController
	Logger     *slog.Logger
}

type Events struct {
	controller IEventsController
	logger     *slog.Logger
}

func NewEvents(dep EventsDependencies) (*Events, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewEvents", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewEvents", "Controller", "nil")
	}

	logger := dep.Logger.With("Handler", "Events")

	return &Events{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Events) NewApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req EventsNewApplication
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.controller.NewApplication(ctx, req)
	if err != nil {
		h.logger.Error(err.Error(), "method", "VerifyIdentity")
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, nil)
}

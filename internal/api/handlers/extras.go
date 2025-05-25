package handlers

import (
	"context"
	"errors"
	"github.com/calyrexx/QuietGrooveBackend/internal/api"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"log/slog"
	"net/http"
)

type IExtrasControllers interface {
	GetAll(ctx context.Context) ([]Extra, error)
	Add(ctx context.Context, extras []Extra) error
	Update(ctx context.Context, extra Extra) error
	Delete(ctx context.Context, extraID int) error
}

type ExtrasDependencies struct {
	Controller IExtrasControllers
	Logger     *slog.Logger
}

type Extras struct {
	controller IExtrasControllers
	logger     *slog.Logger
}

func NewExtras(dep ExtrasDependencies) (*Extras, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewExtras", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewExtras", "Controller", "nil")
	}

	logger := dep.Logger.With("Handler", "Extras")

	return &Extras{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Extras) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	extras, err := h.controller.GetAll(ctx)
	if err != nil {
		h.logger.Error(err.Error(), "method", "GetAll")
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, extras)
}

func (h *Extras) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req []Extra
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.controller.Add(ctx, req); err != nil {
		h.logger.Error(err.Error(), "method", "Add")
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, map[string]string{"message": "extras created"})
}

func (h *Extras) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.URLParamInt(r, "id")
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req Extra
	if err = api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	req.ID = id

	if err = h.controller.Update(ctx, req); err != nil {
		status := http.StatusInternalServerError
		if errors.As(err, &errorspkg.ErrRepoNotFound{}) {
			status = http.StatusNotFound
		}
		h.logger.Error(err.Error(), "method", "Update")
		api.WriteError(w, status, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, map[string]string{"message": "extra updated"})
}

func (h *Extras) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.URLParamInt(r, "id")
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.controller.Delete(ctx, id); err != nil {
		status := http.StatusInternalServerError
		if errors.As(err, &errorspkg.ErrRepoNotFound{}) {
			status = http.StatusNotFound
		}
		h.logger.Error(err.Error(), "method", "Delete")
		api.WriteError(w, status, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, nil)
}

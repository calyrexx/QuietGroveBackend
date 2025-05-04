package handlers

import (
	"context"
	"errors"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/sirupsen/logrus"
	"net/http"
)

type IExtrasControllers interface {
	GetAll(ctx context.Context) ([]Extra, error)
	Add(ctx context.Context, extras Extra) error
	Update(ctx context.Context, extra Extra) error
	Delete(ctx context.Context, extraID int) error
}

type ExtrasDependencies struct {
	Controller IExtrasControllers
	Logger     logrus.FieldLogger
}

type Extras struct {
	controller IExtrasControllers
	logger     logrus.FieldLogger
}

func NewExtras(dep ExtrasDependencies) (*Extras, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewExtras", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewExtras", "Controller", "nil")
	}

	logger := dep.Logger.WithField("Handler", "Extras")

	return &Extras{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Extras) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	extras, err := h.controller.GetAll(ctx)
	if err != nil {
		h.logger.Errorf("get all: %v", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, extras)
}

func (h *Extras) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req Extra
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.controller.Add(ctx, req); err != nil {
		h.logger.Errorf("add: %v", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, map[string]string{"message": "house created"})
}

func (h *Extras) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.URLParamInt(r, "id")
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req Extra
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.controller.Update(ctx, req); err != nil {
		status := http.StatusInternalServerError
		if errors.As(err, &errorspkg.ErrRepoNotFound{}) {
			status = http.StatusNotFound
		}
		h.logger.Errorf("update id=%d: %v", id, err)
		api.WriteError(w, status, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, map[string]string{"message": "house updated"})
}

func (h *Extras) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.URLParamInt(r, "id")
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.controller.Delete(ctx, id); err != nil {
		status := http.StatusInternalServerError
		if errors.As(err, &errorspkg.ErrRepoNotFound{}) {
			status = http.StatusNotFound
		}
		h.logger.Errorf("delete id=%d: %v", id, err)
		api.WriteError(w, status, err)
		return
	}

	api.WriteJSON(w, http.StatusNoContent, nil)
}

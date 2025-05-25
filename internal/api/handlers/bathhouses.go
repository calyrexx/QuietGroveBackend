package handlers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/api"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"log/slog"
	"net/http"
)

type IBathhousesController interface {
	GetAll(ctx context.Context) ([]entities.Bathhouse, error)
	Add(ctx context.Context, bhs []entities.Bathhouse) error
	Update(ctx context.Context, bh entities.Bathhouse) error
	Delete(ctx context.Context, id int) error
}

type BathhousesDependencies struct {
	Controller IBathhousesController
	Logger     *slog.Logger
}

type Bathhouses struct {
	controller IBathhousesController
	logger     *slog.Logger
}

func NewBathhouses(dep BathhousesDependencies) *Bathhouses {
	logger := dep.Logger.With("Handler", "Bathhouses")
	return &Bathhouses{
		controller: dep.Controller,
		logger:     logger,
	}
}

func (h *Bathhouses) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := h.controller.GetAll(ctx)
	if err != nil {
		h.logger.Error("GetAll error", "err", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, list)
}

func (h *Bathhouses) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req []entities.Bathhouse
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.controller.Add(ctx, req); err != nil {
		h.logger.Error("Add error", "err", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, map[string]string{"message": "bathhouses created"})
}

func (h *Bathhouses) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.URLParamInt(r, "id")
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req entities.Bathhouse
	if err = api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	req.ID = id
	if err = h.controller.Update(ctx, req); err != nil {
		h.logger.Error("Update error", "err", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, map[string]string{"message": "bathhouse updated"})
}

func (h *Bathhouses) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.URLParamInt(r, "id")
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.controller.Delete(ctx, id); err != nil {
		h.logger.Error("Delete error", "err", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	api.WriteJSON(w, http.StatusOK, map[string]string{"message": "bathhouse deleted"})
}

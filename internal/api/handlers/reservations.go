package handlers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/api"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/usecases"
	"github.com/gorilla/schema"
	"log/slog"
	"net/http"
)

type IControllers interface {
	CreateReservation(ctx context.Context, req CreateReservation) (entities.Reservation, error)
	GetAvailableHouses(ctx context.Context, req GetAvailableHouses) ([]usecases.GetAvailableHousesResponse, error)
}

type ReservationsDependencies struct {
	Controller IControllers
	Logger     *slog.Logger
}

type Reservations struct {
	controller IControllers
	logger     *slog.Logger
}

func NewReservations(dep ReservationsDependencies) (*Reservations, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewReservations", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewReservations", "Controller", "nil")
	}

	logger := dep.Logger.With("Handler", "Reservations")

	return &Reservations{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Reservations) GetAvailableHouses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := schema.NewDecoder()

	var req GetAvailableHouses
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.controller.GetAvailableHouses(ctx, req)
	if err != nil {
		h.logger.Error(err.Error(), "method", "GetAvailableHouses")
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusOK, result)
}

func (h *Reservations) CreateReservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateReservation
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.controller.CreateReservation(ctx, req)
	if err != nil {
		h.logger.Error(err.Error(), "method", "CreateReservation")
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, result)
}

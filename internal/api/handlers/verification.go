package handlers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/api"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"log/slog"
	"net/http"
)

type IVerificationController interface {
	Generate(ctx context.Context, email, phone string) (string, error)
}

type VerificationDependencies struct {
	Controller IVerificationController
	Logger     *slog.Logger
}

type Verification struct {
	controller IVerificationController
	logger     *slog.Logger
}

func NewVerification(dep VerificationDependencies) (*Verification, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewVerification", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewVerification", "Controller", "nil")
	}

	logger := dep.Logger.With("Handler", "Verification")

	return &Verification{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Verification) VerifyIdentity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req VerifyRequest
	if err := api.ReadJSON(r, &req); err != nil {
		api.WriteError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := h.controller.Generate(ctx, req.Email, req.Phone)
	if err != nil {
		h.logger.Error(err.Error(), "method", "VerifyIdentity")
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, map[string]string{"code": resp})
}

package handlers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/sirupsen/logrus"
	"net/http"
)

type IVerificationController interface {
	Generate(ctx context.Context, email, phone string) (string, error)
}

type VerificationDependencies struct {
	Controller IVerificationController
	Logger     logrus.FieldLogger
}

type Verification struct {
	controller IVerificationController
	logger     logrus.FieldLogger
}

func NewVerification(dep VerificationDependencies) (*Verification, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewVerification", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewVerification", "Controller", "nil")
	}

	logger := dep.Logger.WithField("Handler", "Verification")

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
		h.logger.Errorf("Generate: %v", err)
		api.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	api.WriteJSON(w, http.StatusCreated, map[string]string{"code": resp})
}

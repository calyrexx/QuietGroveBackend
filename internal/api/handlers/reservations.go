package handlers

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/sirupsen/logrus"
	"net/http"
)

type IControllers interface {
}

type ReservationsDependencies struct {
	Controller IControllers
	Logger     logrus.FieldLogger
}

type Reservations struct {
	controller IControllers
	logger     logrus.FieldLogger
}

func NewReservations(dep ReservationsDependencies) (*Reservations, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewReservations", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewReservations", "Controller", "nil")
	}

	logger := dep.Logger.WithField("Handler", "Reservations")

	return &Reservations{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Reservations) BookAHouse(w http.ResponseWriter, r *http.Request) {

}

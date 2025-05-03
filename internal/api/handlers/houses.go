package handlers

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/sirupsen/logrus"
	"net/http"
)

type IHousesControllers interface {
	GetAll(ctx context.Context) ([]entities.House, error)
	Add(ctx context.Context, house entities.House) error
	Update(ctx context.Context, house entities.House) error
	Delete(ctx context.Context, houseID int) error
}

type HousesDependencies struct {
	Controller IHousesControllers
	Logger     logrus.FieldLogger
}

type Houses struct {
	controller IHousesControllers
	logger     logrus.FieldLogger
}

func NewHouses(dep HousesDependencies) (*Houses, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewHouses", "Logger", "nil")
	}
	if dep.Controller == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewHouses", "Controller", "nil")
	}

	logger := dep.Logger.WithField("Handler", "Houses")

	return &Houses{
		controller: dep.Controller,
		logger:     logger,
	}, nil
}

func (h *Houses) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (h *Houses) Add(w http.ResponseWriter, r *http.Request) {

}

func (h *Houses) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *Houses) Delete(w http.ResponseWriter, r *http.Request) {

}

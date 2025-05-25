package controllers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/api/handlers"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/usecases"
	"time"
)

type IReservationsUseCase interface {
	CreateReservation(ctx context.Context, req usecases.CreateReservationRequest) (entities.Reservation, error)
	GetAvailableHouses(ctx context.Context, req entities.GetAvailableHouses) ([]usecases.GetAvailableHousesResponse, error)
}

type ReservationsDependencies struct {
	UseCase IReservationsUseCase
}

type Reservations struct {
	useCase IReservationsUseCase
}

func NewReservations(d *ReservationsDependencies) (*Reservations, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Reservations UseCase", "whole", "nil")
	}
	return &Reservations{
		useCase: d.UseCase,
	}, nil
}

func (c *Reservations) GetAvailableHouses(ctx context.Context, req handlers.GetAvailableHouses) ([]usecases.GetAvailableHousesResponse, error) {
	request, err := c.convertGetAvailableHousesReq(req)
	if err != nil {
		return nil, err
	}

	response, err := c.useCase.GetAvailableHouses(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Reservations) CreateReservation(ctx context.Context, req handlers.CreateReservation) (entities.Reservation, error) {
	request, err := c.convertCreateReservation(req)
	if err != nil {
		return entities.Reservation{}, err
	}

	response, err := c.useCase.CreateReservation(ctx, request)
	if err != nil {
		return entities.Reservation{}, err
	}

	return response, nil
}

func (c *Reservations) convertGetAvailableHousesReq(req handlers.GetAvailableHouses) (entities.GetAvailableHouses, error) {
	in, err := time.Parse(time.DateOnly, req.CheckIn)
	if err != nil {
		return entities.GetAvailableHouses{}, err
	}
	out, err := time.Parse(time.DateOnly, req.CheckOut)
	if err != nil {
		return entities.GetAvailableHouses{}, err
	}
	return entities.GetAvailableHouses{
		CheckIn:     in,
		CheckOut:    out,
		GuestsCount: req.GuestsCount,
	}, nil
}

func (c *Reservations) convertCreateReservation(req handlers.CreateReservation) (usecases.CreateReservationRequest, error) {
	resp := usecases.CreateReservationRequest{
		HouseID:     req.HouseID,
		Guest:       c.convertGuest(req.Guest),
		GuestsCount: req.GuestsCount,
		Extras:      c.convertExtras(req.Extras),
		Bathhouse:   c.convertBathouse(req.Bathhouse),
	}
	cITime, err := time.Parse(time.DateOnly, req.CheckIn)
	if err != nil {
		return resp, err
	}
	cOTime, err := time.Parse(time.DateOnly, req.CheckOut)
	if err != nil {
		return resp, err
	}
	resp.CheckIn = cITime
	resp.CheckOut = cOTime

	return resp, nil
}

func (c *Reservations) convertBathouse(bathhouse []handlers.BathhouseReservation) []entities.BathhouseReservation {
	resp := make([]entities.BathhouseReservation, 0, len(bathhouse))
	for _, b := range bathhouse {
		resp = append(resp, entities.BathhouseReservation{
			TypeID:       b.TypeID,
			Date:         b.Date,
			TimeFrom:     b.TimeFrom,
			TimeTo:       b.TimeTo,
			FillOptionID: b.FillOptionID,
		})
	}
	return resp
}

func (c *Reservations) convertGuest(guest handlers.Guest) entities.Guest {
	return entities.Guest{
		Name:  guest.Name,
		Email: guest.Email,
		Phone: guest.Phone,
	}
}

func (c *Reservations) convertExtras(extras []handlers.ExtraReservation) []entities.ReservationExtra {
	res := make([]entities.ReservationExtra, 0, len(extras))
	for _, e := range extras {
		res = append(res, entities.ReservationExtra{
			ExtraID:  e.ID,
			Quantity: e.Quantity,
			Amount:   e.Amount,
		})
	}
	return res
}

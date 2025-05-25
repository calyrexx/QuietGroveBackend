package usecases

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"github.com/calyrexx/zeroslog"
	"log/slog"
	"time"
)

const reservationStub = "confirmed"

type (
	Notifier interface {
		ReservationCreated(res entities.ReservationCreatedMessage) error
	}

	ReservationDependencies struct {
		ReservationRepo repository.IReservations
		GuestRepo       repository.IGuests
		HouseRepo       repository.IHouses
		PCoefs          []configuration.PriceCoefficient
		Logger          *slog.Logger
		Notifier        Notifier
	}

	Reservation struct {
		reservationRepo repository.IReservations
		guestRepo       repository.IGuests
		houseRepo       repository.IHouses
		pCoefs          []configuration.PriceCoefficient
		logger          *slog.Logger
		notifier        Notifier
	}
)

func NewReservation(d *ReservationDependencies) (*Reservation, error) {
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservation", "whole", "nil")
	}
	if d.ReservationRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservation", "ReservationRepo", "nil")
	}
	if d.GuestRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservation", "GuestRepo", "nil")
	}
	if d.HouseRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservation", "HouseRepo", "nil")
	}
	if d.PCoefs == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservation", "PCoefs", "nil")
	}
	if d.Notifier == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Usecases Reservation", "Notifier", "nil")
	}

	logger := d.Logger.With(zeroslog.UsecaseKey, "Reservation")

	return &Reservation{
		reservationRepo: d.ReservationRepo,
		guestRepo:       d.GuestRepo,
		houseRepo:       d.HouseRepo,
		pCoefs:          d.PCoefs,
		logger:          logger,
		notifier:        d.Notifier,
	}, nil
}

func (u *Reservation) GetAvailableHouses(ctx context.Context, req entities.GetAvailableHouses) ([]GetAvailableHousesResponse, error) {
	availableIDs, err := u.reservationRepo.GetAvailableHouses(ctx, req)
	if err != nil {
		return nil, err
	}

	response := make([]GetAvailableHousesResponse, 0, len(availableIDs))
	for _, id := range availableIDs {
		house, repoErr := u.houseRepo.GetOne(ctx, id)
		if repoErr != nil {
			return nil, repoErr
		}
		nights := int(req.CheckOut.Sub(req.CheckIn).Hours() / 24)
		totalPrice := u.calculateTotalPrice(house.BasePrice, 0, req.CheckIn, req.CheckOut)
		price := totalPrice / nights
		response = append(response, GetAvailableHousesResponse{
			ID:            house.ID,
			Name:          house.Name,
			Description:   house.Description,
			Capacity:      house.Capacity,
			BasePrice:     price,
			TotalPrice:    totalPrice,
			Images:        house.Images,
			CheckInFrom:   house.CheckInFrom,
			CheckOutUntil: house.CheckOutUntil,
		})
	}

	return response, nil
}

func (u *Reservation) CreateReservation(ctx context.Context, req CreateReservationRequest) (entities.Reservation, error) {
	response := entities.Reservation{}

	available, err := u.reservationRepo.CheckAvailability(ctx, entities.CheckAvailability{
		HouseId:  req.HouseID,
		CheckIn:  req.CheckIn,
		CheckOut: req.CheckOut,
	})
	if err != nil {
		return response, err
	}
	if !available {
		return response, errorspkg.NewErrHouseUnavailable(req.HouseID, req.CheckIn, req.CheckOut)
	}

	guest, err := u.guestRepo.FindOrCreate(ctx, req.Guest)
	if err != nil {
		return response, err
	}

	basePrice, err := u.reservationRepo.GetPrice(ctx, req.HouseID, req.Extras)
	if err != nil {
		return response, err
	}

	totalPrice := u.calculateTotalPrice(basePrice.House, basePrice.Extras, req.CheckIn, req.CheckOut)

	reservation := entities.Reservation{
		HouseID:     req.HouseID,
		GuestUUID:   guest.UUID,
		CheckIn:     req.CheckIn,
		CheckOut:    req.CheckOut,
		GuestsCount: req.GuestsCount,
		Status:      reservationStub,
		TotalPrice:  totalPrice,
	}

	if err = u.reservationRepo.Create(ctx, reservation); err != nil {
		return response, err
	}

	go func(res entities.Reservation) {
		house, _ := u.houseRepo.GetOne(context.Background(), res.HouseID)

		reservationMsg := entities.ReservationCreatedMessage{
			House:       house.Name,
			GuestName:   guest.Name,
			GuestPhone:  guest.Phone,
			CheckIn:     res.CheckIn,
			CheckOut:    res.CheckOut,
			GuestsCount: res.GuestsCount,
			TotalPrice:  res.TotalPrice,
		}
		if errSend := u.notifier.ReservationCreated(reservationMsg); errSend != nil {
			u.logger.Error("telegram notify", zeroslog.ErrorKey, err)
		}
	}(reservation)

	return reservation, nil
}

func (u *Reservation) calculateTotalPrice(basePrice, extrasPrice int, checkIn, checkOut time.Time) int {
	total := 0.0
	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	for i := 0; i < nights; i++ {
		currentDate := checkIn.AddDate(0, 0, i)
		coefficient := 1.0

		for _, pc := range u.pCoefs {
			if (currentDate.Equal(pc.Start) || currentDate.After(pc.Start)) &&
				(currentDate.Before(pc.End) || currentDate.Equal(pc.End)) {
				if pc.Rate > coefficient {
					coefficient = pc.Rate
				}
			}
		}

		total += float64(basePrice) * coefficient
	}

	return int(total) + extrasPrice
}

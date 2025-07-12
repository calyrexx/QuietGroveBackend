package usecases

import (
	"context"
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"github.com/calyrexx/zeroslog"
	"log/slog"
	"time"
)

const (
	reservationConfirmed  = "confirmed"
	reservationCheckedIn  = "checked_in"
	reservationCheckedOut = "checked_out"
	barnhouseImg          = "https://res.cloudinary.com/dxmp5yjmb/image/upload/v1747237710/houses1_ebawfo.webp"
	cottageImg            = "https://res.cloudinary.com/dxmp5yjmb/image/upload/v1747237737/houses8_pbv273.jpg"
	glampingImg           = "https://res.cloudinary.com/dxmp5yjmb/image/upload/v1747237765/houses15_djgvjf.webp"
)

type (
	Notifier interface {
		ReservationCreatedForAdmin(res entities.ReservationCreatedMessage) error
		ReservationCreatedForUser(res entities.ReservationCreatedMessage, tgID int64) error
		RemindUser(msg []entities.ReservationReminderNotification) error
	}

	ReservationDependencies struct {
		ReservationRepo repository.IReservations
		GuestRepo       repository.IGuests
		HouseRepo       repository.IHouses
		BathhouseRepo   repository.IBathhouses
		Config          *configuration.Reservations
		Logger          *slog.Logger
		Notifier        Notifier
	}

	Reservation struct {
		reservationRepo repository.IReservations
		guestRepo       repository.IGuests
		houseRepo       repository.IHouses
		bathhouseRepo   repository.IBathhouses
		config          *configuration.Reservations
		logger          *slog.Logger
		notifier        Notifier
	}
)

func NewReservation(d *ReservationDependencies) (*Reservation, error) {
	const method = "usecases.NewReservation"
	if d == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "whole", "nil")
	}
	if d.ReservationRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "ReservationRepo", "nil")
	}
	if d.GuestRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "GuestRepo", "nil")
	}
	if d.HouseRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "HouseRepo", "nil")
	}
	if d.BathhouseRepo == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "BathhouseRepo", "nil")
	}
	if d.Config == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "Config", "nil")
	}
	if d.Notifier == nil {
		return nil, errorspkg.NewErrConstructorDependencies(method, "Notifier", "nil")
	}

	logger := d.Logger.With(zeroslog.UsecaseKey, "Reservation")

	return &Reservation{
		reservationRepo: d.ReservationRepo,
		guestRepo:       d.GuestRepo,
		houseRepo:       d.HouseRepo,
		bathhouseRepo:   d.BathhouseRepo,
		config:          d.Config,
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
		// This is костыль
		bathhouses, _ := u.bathhouseRepo.GetByHouse(ctx, id)
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
			Bathhouses:    u.convertBathhouseToSlots(bathhouses, req.CheckIn, req.CheckOut),
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

	guest, err := u.guestRepo.Get(ctx, req.Guest)
	if err != nil {
		return response, err
	}

	basePrice, err := u.reservationRepo.GetPrice(ctx, req.HouseID, req.Extras, req.Bathhouse)
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
		Status:      reservationConfirmed,
		TotalPrice:  totalPrice,
		Bathhouse:   req.Bathhouse,
	}

	if err = u.reservationRepo.Create(ctx, reservation); err != nil {
		return response, err
	}

	go func(res entities.Reservation, guestTgID int64) {
		house, _ := u.houseRepo.GetOne(context.Background(), res.HouseID)
		bathhouseMsg := make([]entities.BathhouseMessage, 0, len(res.Bathhouse))
		for _, reqBh := range req.Bathhouse {
			bh, _ := u.bathhouseRepo.GetByID(context.Background(), reqBh.TypeID)
			var fillOption *string
			for _, bhFillOptions := range bh.FillOptions {
				if bhFillOptions.ID == reqBh.FillOptionID {
					fillOption = &bhFillOptions.Name
				}
			}
			bathhouseMsg = append(bathhouseMsg, entities.BathhouseMessage{
				Name:       bh.Name,
				Date:       reqBh.Date,
				TimeFrom:   reqBh.TimeFrom,
				TimeTo:     reqBh.TimeTo,
				FillOption: fillOption,
			})
		}

		reservationMsg := entities.ReservationCreatedMessage{
			HouseName:   house.Name,
			GuestName:   guest.Name,
			GuestPhone:  guest.Phone,
			CheckIn:     res.CheckIn,
			CheckOut:    res.CheckOut,
			GuestsCount: res.GuestsCount,
			TotalPrice:  res.TotalPrice,
			Bathhouse:   bathhouseMsg,
		}
		if errSend := u.notifier.ReservationCreatedForAdmin(reservationMsg); errSend != nil {
			u.logger.Error("telegram notify", zeroslog.ErrorKey, err)
		}
		if errSendToUser := u.notifier.ReservationCreatedForUser(reservationMsg, guestTgID); errSendToUser != nil {
			u.logger.Error("telegram user notify", zeroslog.ErrorKey, err)
		}
	}(reservation, guest.TgId)

	return reservation, nil
}

func (u *Reservation) GetByTelegramID(ctx context.Context, userTgID int64) ([]entities.ReservationMessage, error) {
	res, err := u.reservationRepo.GetByTelegramID(ctx, userTgID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *Reservation) GetDetailsByUUID(ctx context.Context, userTgID int64, uuid string) (entities.ReservationMessage, error) {
	res, err := u.reservationRepo.GetDetailsByUUID(ctx, userTgID, uuid)
	if err != nil {
		return entities.ReservationMessage{}, err
	}

	switch res.HouseName {
	case "Барнхаус":
		res.ImageURL = barnhouseImg
	case "Коттедж":
		res.ImageURL = cottageImg
	case "Глэмпинг":
		res.ImageURL = glampingImg
	}

	return res, nil
}

func (u *Reservation) GetForReminder(ctx context.Context) error {
	const method = "GetForReminder"
	u.logger.Info("starting remind users", "method", method)
	timeNow := time.Now()

	reservations, err := u.reservationRepo.GetAllForReminder(ctx)
	if err != nil {
		return err
	}

	notificationDate := timeNow.AddDate(0, 0, u.config.NotificationThreshold)

	reservationsToRemind := make([]entities.ReservationReminderNotification, 0, len(reservations))
	for _, reservation := range reservations {
		if reservation.CheckIn.Year() == notificationDate.Year() &&
			reservation.CheckIn.Month() == notificationDate.Month() &&
			reservation.CheckIn.Day() == notificationDate.Day() {
			reservationsToRemind = append(reservationsToRemind, reservation)
		}
	}

	if err = u.notifier.RemindUser(reservationsToRemind); err != nil {
		return err
	}

	if len(reservationsToRemind) > 0 {
		u.logger.Info(fmt.Sprintf("finished remind users in [%s]", time.Since(timeNow)),
			"method", method, "reservations", len(reservationsToRemind))
	}

	return nil
}

func (u *Reservation) UpdateStatuses(ctx context.Context) error {
	const method = "UpdateStatuses"
	u.logger.Info("starting update statuses", "method", method)
	timeNow := time.Now()

	reservations, err := u.reservationRepo.GetAllConfirmed(ctx)
	if err != nil {
		return err
	}

	reservationsToUpdate := make([]entities.ReservationUpdateStatus, 0, len(reservations))
	for _, reservation := range reservations {
		switch {
		case reservation.CheckIn.Before(timeNow) || reservation.CheckIn.Equal(timeNow):
			if reservation.CheckOut.After(timeNow) {
				reservation.Status = reservationCheckedIn
				reservationsToUpdate = append(reservationsToUpdate, reservation)
			} else if reservation.CheckOut.Before(timeNow) || reservation.CheckOut.Equal(timeNow) {
				reservation.Status = reservationCheckedOut
				reservationsToUpdate = append(reservationsToUpdate, reservation)
			}
		}
	}

	if err = u.reservationRepo.UpdateStatuses(ctx, reservationsToUpdate); err != nil {
		return err
	}

	if len(reservationsToUpdate) > 0 {
		u.logger.Info(fmt.Sprintf("finished update statuses in [%s]", time.Since(timeNow)),
			"method", method, "reservations", len(reservationsToUpdate))
	}

	return nil
}

func (u *Reservation) Cancel(ctx context.Context, userTgID int64, uuid string) error {
	return u.reservationRepo.Cancel(ctx, userTgID, uuid)
}

func (u *Reservation) calculateTotalPrice(basePrice, extrasPrice int, checkIn, checkOut time.Time) int {
	total := 0.0
	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	for i := 0; i < nights; i++ {
		currentDate := checkIn.AddDate(0, 0, i)
		coefficient := 1.0

		for _, pc := range u.config.PriceCoefficients {
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

func (u *Reservation) convertBathhouseToSlots(req []entities.Bathhouse, checkIn, checkOut time.Time) []BathhouseSlots {
	resp := make([]BathhouseSlots, 0, len(req))

	defaultTimeSlot := BathhouseTimeSlots{
		TimeFrom: "10:00",
		TimeTo:   "21:00",
	}

	for _, b := range req {
		days := int(checkOut.Sub(checkIn).Hours() / 24)
		if days == 0 {
			days = 1
		}
		dateSlots := make([]BathhouseDateSlots, 0, days)
		for i := 0; i < days; i++ {
			date := checkIn.AddDate(0, 0, i).Format("2006-01-02")
			dateSlots = append(dateSlots, BathhouseDateSlots{
				Date: date,
				Time: []BathhouseTimeSlots{defaultTimeSlot},
			})
		}
		resp = append(resp, BathhouseSlots{
			TypeID:     b.ID,
			Name:       b.Name,
			Slots:      dateSlots,
			FillOption: u.convertFillOptions(b.FillOptions),
		})
	}
	return resp
}

func (u *Reservation) convertFillOptions(req []entities.BathhouseFillOption) []BathhouseFillOption {
	resp := make([]BathhouseFillOption, 0, len(req))
	for _, b := range req {
		resp = append(resp, BathhouseFillOption{
			ID:          b.ID,
			Name:        b.Name,
			Price:       b.Price,
			Description: b.Description,
		})
	}
	return resp
}

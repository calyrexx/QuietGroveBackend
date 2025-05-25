package postgres

import (
	"context"
	"database/sql"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ReservationsRepo struct {
	pool *pgxpool.Pool
}

func NewReservationsRepo(pool *pgxpool.Pool) *ReservationsRepo {
	return &ReservationsRepo{pool: pool}
}

func (r *ReservationsRepo) GetAvailableHouses(ctx context.Context, req entities.GetAvailableHouses) ([]int, error) {
	const method = "reservationsRepo.CheckAvailability"

	query := `
        SELECT h.id
        FROM houses h
        WHERE h.capacity >= $3
        AND NOT EXISTS (
            SELECT 1 FROM reservations r
            WHERE r.house_id = h.id
            AND r.stay && daterange($1::date, $2::date)
            AND r.status NOT IN ('cancelled', 'checked_out')
        )
        AND NOT EXISTS (
            SELECT 1 FROM blackouts b
            WHERE b.house_id = h.id
            AND b.period && daterange($1::date, $2::date)
        )
    `

	rows, err := r.pool.Query(
		ctx,
		query,
		req.CheckIn,
		req.CheckOut,
		req.GuestsCount,
	)
	if err != nil {
		return nil, errorspkg.NewErrRepoFailed("Query", method, err)
	}
	defer rows.Close()

	var availableHouseIDs []int
	for rows.Next() {
		var houseID int
		if err = rows.Scan(&houseID); err != nil {
			return nil, errorspkg.NewErrRepoFailed("Scan", method, err)
		}
		availableHouseIDs = append(availableHouseIDs, houseID)
	}

	if err = rows.Err(); err != nil {
		return nil, errorspkg.NewErrRepoFailed("rows.Err", method, err)
	}

	return availableHouseIDs, nil
}

func (r *ReservationsRepo) CheckAvailability(ctx context.Context, req entities.CheckAvailability) (bool, error) {
	const method = "reservationsRepo.CheckAvailability"

	query := `
		SELECT NOT EXISTS (
			SELECT 1 FROM reservations
			WHERE house_id = $1 AND stay && daterange($2::date, $3::date)
			UNION ALL
			SELECT 1 FROM blackouts
			WHERE house_id = $1 AND period && daterange($2::date, $3::date)
		)
	`

	var available bool
	err := r.pool.QueryRow(
		ctx,
		query,
		req.HouseId,
		req.CheckIn.Format(time.DateOnly),
		req.CheckOut.Format(time.DateOnly),
	).Scan(&available)

	if err != nil {
		return false, errorspkg.NewErrRepoFailed("QueryRow", method, err)
	}

	return available, nil
}

func (r *ReservationsRepo) GetPrice(ctx context.Context, houseID int, extras []entities.ReservationExtra, bathhouses []entities.BathhouseReservation) (entities.GetPrice, error) {
	const method = "reservationsRepo.GetPrice"

	response := entities.GetPrice{}
	housePrice, extrasPrice, bathhousePrice := 0, 0, 0
	housePriceQuery := `
		SELECT base_price FROM houses 
			WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, housePriceQuery, houseID).Scan(&housePrice)
	if err != nil {
		return response, errorspkg.NewErrRepoFailed("QueryRow", method, err)
	}

	response.House = housePrice

	if len(extras) > 0 {
		extraIds := make([]int, 0, len(extras))
		for _, e := range extras {
			extraIds = append(extraIds, e.ExtraID)
		}

		extrasPriceQuery := `
			SELECT price FROM extras 
            	WHERE id = ANY($1)
		`
		rows, extErr := r.pool.Query(ctx, extrasPriceQuery, extraIds)
		if extErr != nil {
			return response, errorspkg.NewErrRepoFailed("Query", method, err)
		}
		defer rows.Close()

		for rows.Next() {
			var price int
			if err = rows.Scan(&price); err != nil {
				return response, errorspkg.NewErrRepoFailed("Scan", method, err)
			}
			extrasPrice += price
		}

		response.Extras = extrasPrice
	}

	if len(bathhouses) > 0 {
		bathhouseIds := make([]int, 0, len(bathhouses))
		for _, b := range bathhouses {
			bathhouseIds = append(bathhouseIds, b.TypeID)
		}

		bathhousesPriceQuery := `
			SELECT default_price FROM sauna_types 
				WHERE id = ANY($1)
		`
		rows, extErr := r.pool.Query(ctx, bathhousesPriceQuery, bathhouseIds)
		if extErr != nil {
			return response, errorspkg.NewErrRepoFailed("Query", method, err)
		}
		defer rows.Close()

		for rows.Next() {
			var price int
			if err = rows.Scan(&price); err != nil {
				return response, errorspkg.NewErrRepoFailed("Scan", method, err)
			}
			bathhousePrice += price
		}

		response.Bathhouse = bathhousePrice
	}

	return response, nil
}

func (r *ReservationsRepo) Create(ctx context.Context, reservation entities.Reservation) error {
	const method = "reservationsRepo.Create"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errorspkg.NewErrRepoFailed("BeginTx", method, err)
	}

	var resUUID uuid.UUID
	queryReservation := `
		INSERT INTO reservations (
			uuid, house_id, guest_uuid, stay, guests_count, status, total_price
		) VALUES (
			$1, $2, $3, daterange($4::date, $5::date), $6, $7, $8
		)
		RETURNING uuid
	`
	resUUID = uuid.New()
	_, err = tx.Exec(ctx, queryReservation,
		resUUID,
		reservation.HouseID,
		reservation.GuestUUID,
		reservation.CheckIn.Format(time.DateOnly),
		reservation.CheckOut.Format(time.DateOnly),
		reservation.GuestsCount,
		reservation.Status,
		reservation.TotalPrice,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return errorspkg.NewErrRepoFailed("Exec Insert Reservation", method, err)
	}

	if len(reservation.Bathhouse) > 0 {
		queryBath := `
			INSERT INTO bathhouse_reservations (
				reservation_uuid, type_id, date, time_from, time_to, fill_option_id
			) VALUES ($1, $2, $3::date, $4::time, $5::time, $6)
		`
		for _, b := range reservation.Bathhouse {
			_, err = tx.Exec(ctx, queryBath,
				resUUID,
				b.TypeID,
				b.Date,
				b.TimeFrom,
				b.TimeTo,
				sql.NullInt64{Int64: int64(b.FillOptionID), Valid: b.FillOptionID != 0},
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				return errorspkg.NewErrRepoFailed("Exec Insert Bathhouse", method, err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return errorspkg.NewErrRepoFailed("Commit", method, err)
	}

	return nil
}

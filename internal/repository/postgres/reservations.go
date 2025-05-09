package postgres

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationsRepo struct {
	pool *pgxpool.Pool
}

func NewReservationsRepo(pool *pgxpool.Pool) *ReservationsRepo {
	return &ReservationsRepo{pool: pool}
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
		req.CheckIn.Format(repository.DateFormat),
		req.CheckOut.Format(repository.DateFormat),
	).Scan(&available)

	if err != nil {
		return false, errorspkg.NewErrRepoFailed("QueryRow", method, err)
	}

	return available, nil
}

func (r *ReservationsRepo) GetPrice(ctx context.Context, houseID int, extras []entities.ReservationExtra) (entities.GetPrice, error) {
	const method = "reservationsRepo.GetPrice"

	response := entities.GetPrice{}
	housePrice, extrasPrice := 0, 0
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

		extrasPriceQuery := `SELECT price FROM extras WHERE id = ANY($1)`
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

	return response, nil
}

func (r *ReservationsRepo) Create(ctx context.Context, reservation entities.Reservation) error {
	const method = "reservationsRepo.Create"

	query := `
		INSERT INTO reservations (
			uuid, house_id, guest_uuid, stay, guests_count, status, total_price
		) VALUES (
			$1, $2, $3, daterange($4::date, $5::date), $6, $7, $8
		)
	`

	_, err := r.pool.Exec(ctx, query,
		uuid.New(),
		reservation.HouseID,
		reservation.GuestUUID,
		reservation.CheckIn.Format(repository.DateFormat),
		reservation.CheckOut.Format(repository.DateFormat),
		reservation.GuestsCount,
		reservation.Status,
		reservation.TotalPrice,
	)

	if err != nil {
		return errorspkg.NewErrRepoFailed("Exec", method, err)
	}
	return nil
}

package postgres

import (
	"context"
	"errors"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

type HousesRepo struct {
	pool *pgxpool.Pool
}

func NewHousesRepo(pool *pgxpool.Pool) *HousesRepo {
	return &HousesRepo{pool: pool}
}

func (r *HousesRepo) GetAll(ctx context.Context) ([]entities.House, error) {
	const method = "housesRepo.GetAll"
	rows, err := r.pool.Query(ctx, `
		SELECT 
			id,
			name,
			description,
			
			capacity,
			base_price,
			
			images,
			
			check_in_from,
			check_out_until
		FROM houses
	`)
	if err != nil {
		return nil, errorspkg.NewErrRepoFailed("pool.Query", method, err)
	}
	defer rows.Close()

	var results []entities.House
	for rows.Next() {
		var house entities.House
		if err = rows.Scan(
			&house.ID,
			&house.Name,
			&house.Description,
			&house.Capacity,
			&house.BasePrice,
			&house.Images,
			&house.CheckInFrom,
			&house.CheckOutUntil,
		); err != nil {
			return nil, errorspkg.NewErrRepoFailed("rows.Scan", method, err)
		}
		results = append(results, house)
	}
	if err = rows.Err(); err != nil {
		return nil, errorspkg.NewErrRepoFailed("rows.Err", method, err)
	}

	return results, nil
}

func (r *HousesRepo) GetOne(ctx context.Context, id int) (entities.House, error) {
	const method = "housesRepo.GetOne"

	query := `
        SELECT 
            id,
            name,
            description,
            
            capacity,
            base_price,
            
            images,
            
            check_in_from,
            check_out_until
        FROM houses
        WHERE id = $1
    `

	var house entities.House

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&house.ID,
		&house.Name,
		&house.Description,
		&house.Capacity,
		&house.BasePrice,
		&house.Images,
		&house.CheckInFrom,
		&house.CheckOutUntil,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.House{}, errorspkg.NewErrRepoNotFound("house", strconv.Itoa(id), method)
		}
		return entities.House{}, errorspkg.NewErrRepoFailed("QueryRow", method, err)
	}

	return house, nil
}

func (r *HousesRepo) Add(ctx context.Context, house entities.House) error {
	const method = "housesRepo.Add"
	query := `
		INSERT INTO houses (
		    id,
			name,
		    description,
		                    
			capacity,
			base_price,
		                    
			images,
		                    
			check_in_from,
			check_out_until
		)
		VALUES (
			$1, $2, $3, 
		    $4, $5,
		    $6, 
		    $7, $8
		)
	`

	if _, err := r.pool.Exec(ctx, query,
		house.ID,
		house.Name,
		house.Description,
		house.Capacity,
		house.BasePrice,
		house.Images,
		house.CheckInFrom,
		house.CheckOutUntil,
	); err != nil {
		return errorspkg.NewErrRepoFailed("pool.Exec", method, err)
	}

	return nil
}

func (r *HousesRepo) Update(ctx context.Context, house entities.House) error {
	const method = "housesRepo.Update"
	query := `
		UPDATE houses
		SET
			name            = $1,
			capacity        = $2,
			base_price      = $3,
			description     = $4,
		    images          = $5,
		    check_in_from   = $6,
		    check_out_until = $7,
		    updated_at      = now()
		WHERE id = $8
	`

	rows, err := r.pool.Exec(ctx, query,
		house.Name,
		house.Capacity,
		house.BasePrice,
		house.Description,
		house.Images,
		house.CheckInFrom,
		house.CheckOutUntil,
		house.ID,
	)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Exec", method, err)
	}
	if rows.RowsAffected() == 0 {
		return errorspkg.NewErrRepoNotFound("house", strconv.Itoa(house.ID), method)
	}
	return nil
}

func (r *HousesRepo) Delete(ctx context.Context, id int) error {
	const method = "housesRepo.Delete"
	query := `
		DELETE FROM houses WHERE id = $1
	`

	rows, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Exec", method, err)
	}
	if rows.RowsAffected() == 0 {
		return errorspkg.NewErrRepoNotFound("house", strconv.Itoa(id), method)
	}
	return nil
}

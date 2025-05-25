package postgres

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

type ExtrasRepo struct {
	pool *pgxpool.Pool
}

func NewExtrasRepo(pool *pgxpool.Pool) *ExtrasRepo {
	return &ExtrasRepo{pool: pool}
}

func (r *ExtrasRepo) GetAll(ctx context.Context) ([]entities.Extra, error) {
	const method = "extrasRepo.GetAll"
	rows, err := r.pool.Query(ctx, `
		SELECT 
			id,
			name,
			short_text,
			description,
			price,
			images
		FROM extras
		ORDER BY id
	`)
	if err != nil {
		return nil, errorspkg.NewErrRepoFailed("pool.Query", method, err)
	}
	defer rows.Close()

	var extras []entities.Extra
	for rows.Next() {
		var e entities.Extra
		if err = rows.Scan(
			&e.ID,
			&e.Name,
			&e.Text,
			&e.Description,
			&e.BasePrice,
			&e.Images,
		); err != nil {
			return nil, errorspkg.NewErrRepoFailed("rows.Scan", method, err)
		}
		extras = append(extras, e)
	}
	if err = rows.Err(); err != nil {
		return nil, errorspkg.NewErrRepoFailed("rows.Err", method, err)
	}
	return extras, nil
}

func (r *ExtrasRepo) Add(ctx context.Context, extras []entities.Extra) error {
	const method = "extrasRepo.Add"

	query := `
		INSERT INTO extras (
			id,
			name,
		    short_text,
			description,
			price,
			images
		)
		VALUES (
			$1, 
		    $2,
		    $3, 
		    $4,
		    $5,
		    $6
		)
	`

	batch := &pgx.Batch{}
	for _, extra := range extras {
		batch.Queue(query,
			extra.ID,
			extra.Name,
			extra.Text,
			extra.Description,
			extra.BasePrice,
			extra.Images,
		)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(extras); i++ {
		if _, err := br.Exec(); err != nil {
			return errorspkg.NewErrRepoFailed("batch.Exec", method, err)
		}
	}

	return nil
}

func (r *ExtrasRepo) Update(ctx context.Context, extra entities.Extra) error {
	const method = "extrasRepo.Update"
	query := `
		UPDATE extras
		SET
		    name        = $1,
		    description = $2,
		    price       = $3,
		    images      = $4,
		    short_text   = $5,
		    updated_at  = now()
		WHERE id = $6
	`
	rows, err := r.pool.Exec(ctx, query,
		extra.Name,
		extra.Description,
		extra.BasePrice,
		extra.Images,
		extra.Text,
		extra.ID,
	)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Exec", method, err)
	}
	if rows.RowsAffected() == 0 {
		return errorspkg.NewErrRepoNotFound("extra", strconv.Itoa(extra.ID), method)
	}
	return nil
}

func (r *ExtrasRepo) Delete(ctx context.Context, id int) error {
	const method = "extrasRepo.Delete"

	tag, err := r.pool.Exec(ctx, `DELETE FROM extras WHERE id = $1`, id)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Exec", method, err)
	}
	if tag.RowsAffected() == 0 {
		return errorspkg.NewErrRepoNotFound("extra", strconv.Itoa(id), method)
	}

	return nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuestsRepo struct {
	pool *pgxpool.Pool
}

func NewGuestsRepo(pool *pgxpool.Pool) *GuestsRepo {
	return &GuestsRepo{pool: pool}
}

func (r *GuestsRepo) Create(ctx context.Context, guest entities.Guest) error {
	const method = "guestsRepo.Create"

	query := `
		INSERT INTO guests (uuid, name, email, phone, tg_user_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		uuid.New(),
		guest.Name,
		guest.Email,
		guest.Phone,
		guest.TgID,
	)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Exec", method, err)
	}

	return nil
}

func (r *GuestsRepo) Get(ctx context.Context, req entities.Guest) (repository.Guest, error) {
	const method = "guestsRepo.Get"

	query := `
		SELECT uuid, name, email, phone, tg_user_id 
		FROM guests 
		WHERE email = $1 AND phone = $2 AND name = $3
	`

	var guest repository.Guest
	err := r.pool.QueryRow(ctx, query, req.Email, req.Phone, req.Name).Scan(
		&guest.UUID,
		&guest.Name,
		&guest.Email,
		&guest.Phone,
		&guest.TgId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return guest, errorspkg.NewErrRepoNotFound("guest", fmt.Sprintf("%s %s %s", req.Name, req.Phone, req.Email), method)
		}
		return guest, errorspkg.NewErrRepoFailed("QueryRow", method, err)
	}

	return guest, nil
}

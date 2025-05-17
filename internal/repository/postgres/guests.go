package postgres

import (
	"context"
	"errors"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository"
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

func (r *GuestsRepo) FindOrCreate(ctx context.Context, guest entities.Guest) (repository.Guest, error) {
	const method = "guestsRepo.FindOrCreate"

	existing, err := r.findByEmail(ctx, guest)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return repository.Guest{}, err
	}

	query := `
		INSERT INTO guests (uuid, name, email, phone)
		VALUES ($1, $2, $3, $4)
	`

	newGuest := repository.Guest{
		UUID:  uuid.New(),
		Name:  guest.Name,
		Email: guest.Email,
		Phone: guest.Phone,
	}

	_, err = r.pool.Exec(ctx, query,
		newGuest.UUID,
		newGuest.Name,
		newGuest.Email,
		newGuest.Phone,
	)
	if err != nil {
		return repository.Guest{}, errorspkg.NewErrRepoFailed("Exec", method, err)
	}

	return newGuest, nil
}

func (r *GuestsRepo) findByEmail(ctx context.Context, req entities.Guest) (repository.Guest, error) {
	query := `
		SELECT uuid, name, email, phone 
		FROM guests 
		WHERE email = $1 OR phone = $2
	`

	var guest repository.Guest
	err := r.pool.QueryRow(ctx, query, req.Email, req.Phone).Scan(
		&guest.UUID,
		&guest.Name,
		&guest.Email,
		&guest.Phone,
	)
	if err != nil {
		return guest, err
	}

	return guest, nil
}

package postgres

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationRepo struct {
	pool *pgxpool.Pool
}

func NewVerificationRepo(pool *pgxpool.Pool) *VerificationRepo {
	return &VerificationRepo{
		pool: pool,
	}
}

func (r *VerificationRepo) Create(ctx context.Context, v entities.Verification) error {
	query := `
		INSERT INTO verifications (
		    uuid,
		    name,
			code,
		    email,
		    phone,
		    status,
		    expires_at
		)
        VALUES (
        	$1,
        	$2,
            $3,
            $4,
            $5,
            $6,
            $7
        )`

	_, err := r.pool.Exec(ctx, query, uuid.New(), v.Name, v.Code, v.Email, v.Phone, v.Status, v.ExpiresAt)

	return err
}

func (r *VerificationRepo) GetByCode(ctx context.Context, code string) (entities.Verification, error) {
	query := `
		SELECT 
			uuid,
		    code,
		    email,
		    phone,
		    tg_user_id,
		    status,
		    created_at,
		    verified_at,
		    expires_at
        FROM verifications WHERE code=$1`

	var v entities.Verification
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&v.ID,
		&v.Code,
		&v.Email,
		&v.Phone,
		&v.TgUserID,
		&v.Status,
		&v.CreatedAt,
		&v.VerifiedAt,
		&v.ExpiresAt,
	)

	return v, err
}

func (r *VerificationRepo) Approve(ctx context.Context, id string, tgUserID int64) error {
	query := `
		UPDATE verifications
        SET status=$1, tg_user_id=$2, verified_at=now()
        WHERE uuid=$3`

	_, err := r.pool.Exec(ctx, query,
		entities.VerifApproved,
		tgUserID,
		id,
	)

	return err
}

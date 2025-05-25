package postgres

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
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

func (r *VerificationRepo) Create(ctx context.Context, v entities.Verification) (int64, error) {
	var id int64

	query := `
		INSERT INTO verifications (
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
            $5
        ) 
        RETURNING id`

	err := r.pool.QueryRow(ctx, query, v.Code, v.Email, v.Phone, v.Status, v.ExpiresAt).Scan(&id)

	return id, err
}

func (r *VerificationRepo) GetByCode(ctx context.Context, code string) (entities.Verification, error) {
	query := `
		SELECT 
			id,
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

func (r *VerificationRepo) Approve(ctx context.Context, id int64, tgUserID int64) error {
	query := `
		UPDATE verifications
        SET status=$1, tg_user_id=$2, verified_at=now()
        WHERE id=$3`

	_, err := r.pool.Exec(ctx, query,
		entities.VerifApproved,
		tgUserID,
		id,
	)

	return err
}

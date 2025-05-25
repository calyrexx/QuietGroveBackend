package usecases

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/repository"
	"math/big"
	"time"
)

type VerificationDependencies struct {
	Repo repository.IVerification
	TTL  time.Duration
}

type Verification struct {
	repo repository.IVerification
	ttl  time.Duration
}

func NewVerification(d *VerificationDependencies) (*Verification, error) {
	if d.Repo == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewVerification", "Repo", "nil")
	}
	if d.TTL == 0 {
		return nil, errorspkg.NewErrConstructorDependencies("NewVerification", "TTL", "0")
	}

	return &Verification{
		repo: d.Repo,
		ttl:  d.TTL,
	}, nil
}

func (s *Verification) Generate(ctx context.Context, email, phone, name string) (string, error) {
	code := sixDigits()
	exp := time.Now().Add(s.ttl)

	_, err := s.repo.Create(ctx, entities.Verification{
		Code:      code,
		Email:     email,
		Phone:     phone,
		Name:      name,
		Status:    entities.VerifPending,
		ExpiresAt: exp,
	})
	return code, err
}

func (s *Verification) Approve(ctx context.Context, code string, tgID int64) error {
	v, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return err
	}

	if v.Status != entities.VerifPending || time.Now().After(v.ExpiresAt) {
		return errorspkg.ErrInvalidVerificationCode
	}

	return s.repo.Approve(ctx, v.ID, tgID)
}

func sixDigits() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900_000))
	if err != nil {
		return ""
	}
	code := n.Int64() + 100_000
	return fmt.Sprintf("%06d", code)
}

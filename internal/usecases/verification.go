package usecases

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/repository"
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

func (s *Verification) Generate(ctx context.Context, email, phone string) (string, error) {
	code := sixDigits()
	exp := time.Now().Add(s.ttl)

	_, err := s.repo.Create(ctx, entities.Verification{
		Code:      code,
		Email:     email,
		Phone:     phone,
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
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2])%1e6)
}

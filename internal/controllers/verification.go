package controllers

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
)

type IVerificationUseCase interface {
	Generate(ctx context.Context, email, phone, name string) (string, error)
}

type VerificationDependencies struct {
	UseCase IVerificationUseCase
}

type Verification struct {
	useCase IVerificationUseCase
}

func NewVerification(d *VerificationDependencies) (*Verification, error) {
	if d.UseCase == nil {
		return nil, errorspkg.NewErrConstructorDependencies("Verification UseCase", "whole", "nil")
	}

	return &Verification{
		useCase: d.UseCase,
	}, nil
}

func (c *Verification) Generate(ctx context.Context, email, phone, name string) (string, error) {
	return c.useCase.Generate(ctx, email, phone, name)
}

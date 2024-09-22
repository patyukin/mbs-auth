package usecase

import (
	"context"
	"github.com/patyukin/mbs-auth/internal/db"
)

type Registry interface {
	GetRepo() db.RepositoryInterface
	ReadCommitted(ctx context.Context, f db.Handler) error
}

type UseCase struct {
	registry Registry
}

func New(registry Registry) *UseCase {
	return &UseCase{
		registry: registry,
	}
}

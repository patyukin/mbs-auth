package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"github.com/rs/zerolog/log"
)

type RepositoryInterface interface {
	InsertIntoUsers(ctx context.Context, in model.User) (uuid.UUID, error)
	InsertIntoProfiles(ctx context.Context, in model.Profile) (uuid.UUID, error)
	InsertIntoTelegramUsers(ctx context.Context, in model.TelegramUser) (uuid.UUID, error)
	SelectUsersWithTokensCount(ctx context.Context) (int32, error)
	SelectUsersWithTokens(ctx context.Context, limit int32, page int32) ([]*authpb.UserGUWR, error)
	SelectUsersWithProfilesCount(ctx context.Context) (int32, error)
	SelectUsersWithProfiles(ctx context.Context, limit int32, page int32) ([]model.UserWithProfile, error)
}

type Registry struct {
	db *sql.DB
}

func (registry *Registry) GetRepo() RepositoryInterface {
	return &Repository{
		db: registry.db,
	}
}

type Handler func(ctx context.Context, repo RepositoryInterface) error

func New(db *sql.DB) *Registry {
	return &Registry{db: db}
}

func (registry *Registry) ReadCommitted(ctx context.Context, f Handler) error {
	tx, err := registry.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed registry.db.BeginTx: %w", err)
	}

	defer func() {
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Error().Msgf("failed to rollback transaction: %v", errRollback)
			}
		}
	}()

	repo := &Repository{db: tx}

	if err = f(ctx, repo); err != nil {
		return fmt.Errorf("failed to execute handler: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (registry *Registry) Close() error {
	return registry.db.Close()
}

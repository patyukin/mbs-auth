package usecase

import (
	"context"
	"fmt"

	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) RemoveNotRegisteredUsers(ctx context.Context) error {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		userIDs, err := repo.SelectNotRegisteredUsers(ctx)
		if err != nil {
			return fmt.Errorf("failed repo.SelectNotRegisteredUsers: %w", err)
		}

		if len(userIDs) == 0 {
			log.Info().Msg("no not registered users")
			return nil
		}

		err = repo.RemoveTelegramUsersByUserIDs(ctx, userIDs)
		if err != nil {
			return fmt.Errorf("failed removing not registered users in repo.RemoveTelegramUsersByUserIDs: %w", err)
		}

		err = repo.RemoveProfilesByUserIDs(ctx, userIDs)
		if err != nil {
			return fmt.Errorf("failed removing not registered users in repo.RemoveProfilesByUserIDs: %w", err)
		}

		err = repo.RemoveUsersByIDs(ctx, userIDs)
		if err != nil {
			return fmt.Errorf("failed removing not registered users in repo.RemoveUsersByIDs: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed removing not registered users: %w", err)
	}

	log.Info().Msgf("Removed not registered users")
	return nil
}

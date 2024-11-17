package usecase

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) CleanExpiredTokens(ctx context.Context) error {
	err := u.registry.GetRepo().CleanExpiredTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed cleaning tokens: %w", err)
	}

	log.Info().Msgf("Cleaned expired tokens")
	return nil
}

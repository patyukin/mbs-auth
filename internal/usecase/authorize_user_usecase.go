package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) AuthorizeUserV1UseCase(ctx context.Context, in *authpb.AuthorizeUserRequest) (*authpb.AuthorizeUserResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		exists, err := repo.SelectExistsRowByRoleUserIDAndRoutePath(ctx, in)
		log.Debug().Msgf("exists: %v", exists)
		if err != nil {
			return fmt.Errorf("failed repo.SelectExistsRowByRoleUserIDAndRoutePath: %w", err)
		}

		if !exists {
			return fmt.Errorf("no access, %w", errs.ErrUserNotFound)
		}

		return nil
	},
	)
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &authpb.AuthorizeUserResponse{Message: "Authorized"}, nil
}

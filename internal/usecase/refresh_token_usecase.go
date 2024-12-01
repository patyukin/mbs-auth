package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/db"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (u *UseCase) RefreshTokenV1UseCase(ctx context.Context, in *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	var token string

	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		userID, role, err := repo.SelectByID(ctx, in.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed repo.SelectByID: %w", err)
		}

		token, err = u.generateJWT(userID, role)
		if err != nil {
			return fmt.Errorf("failed u.generateJWT: %w", err)
		}

		return nil
	},
	)
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &authpb.RefreshTokenResponse{AccessToken: token}, nil
}

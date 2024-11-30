package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (u *UseCase) SignInConfirmationV1UseCase(ctx context.Context, in *authpb.SignInConfirmationRequest) (*authpb.SignInConfirmationResponse, error) {
	userID, err := u.chr.Get2FACode(ctx, in.Code)
	if err != nil {
		return nil, fmt.Errorf("failed u.chr.Get2FACode: %w", err)
	}

	err = u.chr.Delete2FACode(ctx, in.Code)
	if err != nil {
		return nil, fmt.Errorf("failed u.chr.Delete2FACode: %w", err)
	}

	var user model.User
	var token string
	var refreshToken string

	err = u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			user, err = repo.SelectUserByUUID(ctx, userID)
			if err != nil {
				return fmt.Errorf("failed repo.SelectUserByUUID: %w", err)
			}

			token, err = u.generateJWT(user.UUID.String())
			if err != nil {
				return fmt.Errorf("failed u.generateJWT: %w", err)
			}

			refreshToken, err = repo.UpsertToken(ctx, user.UUID)
			if err != nil {
				return fmt.Errorf("failed repo.UpsertToken: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to u.registry.ReadCommitted: %w", err)
	}

	return &authpb.SignInConfirmationResponse{AccessToken: token, RefreshToken: refreshToken}, nil
}

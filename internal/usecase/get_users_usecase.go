package usecase

import (
	"context"
	"fmt"

	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (u *UseCase) GetUsersV1UseCase(ctx context.Context, in *authpb.GetUsersRequest) (*authpb.GetUsersResponse, error) {
	var response *authpb.GetUsersResponse
	err := u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			totalPages, err := repo.SelectUsersWithProfilesCount(ctx)
			if err != nil {
				return fmt.Errorf("failed repo.SelectUsersWithTokensCount: %w", err)
			}

			users, err := repo.SelectUsersWithProfiles(ctx, in.GetLimit(), in.GetPage())
			if err != nil {
				return fmt.Errorf("failed repo.SelectUsersWithTokens: %w", err)
			}

			usersForResponse := model.ToProtoUserInfo(users)
			response = &authpb.GetUsersResponse{
				Users: usersForResponse,
				Total: totalPages,
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return response, nil
}

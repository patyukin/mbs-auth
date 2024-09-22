package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
)

func (u *UseCase) GetUsersWithProfiles(ctx context.Context, in *authpb.GetUsersWithProfilesRequest) (*authpb.GetUsersWithProfilesResponse, error) {
	if in.Page < 1 {
		in.Page = 1
	}

	if in.Limit < 1 {
		in.Limit = 30
	}

	var response *authpb.GetUsersWithProfilesResponse

	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo db.RepositoryInterface) error {
		totalPages, err := repo.SelectUsersWithProfilesCount(ctx)
		if err != nil {
			return fmt.Errorf("failed repo.SelectUsersWithTokensCount: %w", err)
		}

		users, err := repo.SelectUsersWithProfiles(ctx, in.Limit, in.Page)
		if err != nil {
			return fmt.Errorf("failed repo.SelectUsersWithTokens: %w", err)
		}

		usersForResponse := model.UsersWithProfilesResponseFromUserWithProfile(users)
		response = &authpb.GetUsersWithProfilesResponse{
			Users: usersForResponse,
			Total: totalPages,
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return response, nil
}

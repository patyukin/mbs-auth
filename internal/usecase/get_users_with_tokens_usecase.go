package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/db"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
)

func (u *UseCase) GetUsersWithTokens(ctx context.Context, in *authpb.GetUsersWithTokensRequest) (*authpb.GetUsersWithTokensResponse, error) {
	if in.Page < 1 {
		in.Page = 1
	}

	if in.Limit < 1 {
		in.Limit = 30
	}

	var response *authpb.GetUsersWithTokensResponse

	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo db.RepositoryInterface) error {
		total, err := repo.SelectUsersWithProfilesCount(ctx)
		if err != nil {
			return fmt.Errorf("failed repo.SelectUsersWithProfilesCount: %w", err)
		}

		users, err := repo.SelectUsersWithTokens(ctx, in.Limit, in.Page)
		if err != nil {
			return fmt.Errorf("failed repo.SelectUsersWithProfiles: %w", err)
		}

		response = &authpb.GetUsersWithTokensResponse{
			Users: users,
			Total: total,
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return response, nil
}

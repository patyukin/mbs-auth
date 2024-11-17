package usecase

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (u *UseCase) AddUserRole(ctx context.Context, in *authpb.AddUserRoleRequest) (*authpb.AddUserRoleResponse, error) {
	err := u.registry.GetRepo().InsertIntoUsersRoles(ctx, in.UserId, in.RoleId)
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.GetRepo().InsertIntoUsersRoles: %w", err)
	}

	return &authpb.AddUserRoleResponse{}, nil
}

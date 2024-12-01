package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) GetUsers(ctx context.Context, in *desc.GetUsersRequest) (*desc.GetUsersResponse, error) {
	users, err := s.uc.GetUsersV1UseCase(ctx, in)
	if err != nil {
		return &desc.GetUsersResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.GetUsersV1UseCase: %w", err)),
		}, nil
	}

	if users.Error != nil {
		return &desc.GetUsersResponse{Error: users.Error}, nil
	}

	return users, nil
}

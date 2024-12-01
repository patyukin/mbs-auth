package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) GetUserByID(ctx context.Context, in *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error) {
	userInfo, err := s.uc.GetUserByIDUseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.GetUserByIDUseCase: %w", err)
	}

	return userInfo, nil
}

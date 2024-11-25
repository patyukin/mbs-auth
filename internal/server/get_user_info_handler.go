package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) GetUserInfo(ctx context.Context, in *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error) {
	userInfo, err := s.uc.GetUserInfoUseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.GetUserInfoUseCase: %w", err)
	}

	return userInfo, nil
}

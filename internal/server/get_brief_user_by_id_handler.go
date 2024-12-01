package server

import (
	"context"
	"fmt"

	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) GetBriefUserByID(ctx context.Context, in *authpb.GetBriefUserByIDRequest) (*authpb.GetBriefUserByIDResponse, error) {
	response, err := s.uc.GetBriefUserById(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.GetBriefUserByIDUseCase: %w", err)
	}

	return response, nil
}

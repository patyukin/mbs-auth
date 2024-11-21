package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) Authorize(ctx context.Context, in *authpb.AuthorizeRequest) (*authpb.AuthorizeResponse, error) {
	response, err := s.uc.Authorize(ctx, in)
	if err != nil || response.Error != nil {
		return nil, fmt.Errorf("failed s.uc.Authorize: %w", err)
	}

	return &authpb.AuthorizeResponse{}, nil
}

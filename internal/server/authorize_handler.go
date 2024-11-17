package server

import (
	"context"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

func (s *Server) Authorize(ctx context.Context, in *authpb.AuthorizeRequest) (*authpb.AuthorizeResponse, error) {
	log.Debug().Msgf("Authorize: %v", in)

	return &authpb.AuthorizeResponse{}, nil
}

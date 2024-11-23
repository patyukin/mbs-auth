package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

type UseCase interface {
	SignUpV1UseCase(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error)
	SignInV1UseCase(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error)
	SignInVerifyV1UseCase(ctx context.Context, in *authpb.SignInConfirmationRequest) (*authpb.SignInConfirmationResponse, error)
	GetUsersV1UseCase(ctx context.Context, in *authpb.GetUsersRequest) (*authpb.GetUsersResponse, error)
	AddUserRoleV1UseCase(ctx context.Context, in *authpb.AddUserRoleRequest) (*authpb.AddUserRoleResponse, error)
	AuthorizeUserV1UseCase(ctx context.Context, in *authpb.AuthorizeUserRequest) (*authpb.AuthorizeUserResponse, error)
	RefreshTokenV1UseCase(ctx context.Context, in *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error)
}

type Server struct {
	authpb.UnimplementedAuthServiceServer
	uc UseCase
}

func New(uc UseCase) *Server {
	return &Server{
		uc: uc,
	}
}

func (s *Server) GetUsersWithTokens(ctx context.Context, in *authpb.GetUsersWithTokensRequest) (*authpb.GetUsersWithTokensResponse, error) {
	response, err := s.uc.GetUsersWithTokens(ctx, in)
	if err != nil {
		log.Error().Msgf("failed uc.GetUsersWithTokens: %v", err)
		return nil, fmt.Errorf("failed GetUsersWithTokens")
	}

	return response, nil
}

func (s *Server) GetUsersWithProfiles(ctx context.Context, in *authpb.GetUsersWithProfilesRequest) (*authpb.GetUsersWithProfilesResponse, error) {
	response, err := s.uc.GetUsersWithProfiles(ctx, in)
	if err != nil {
		log.Error().Msgf("failed uc.GetUsersWithProfiles: %v", err)
		return nil, fmt.Errorf("failed GetUsersWithProfiles")
	}

	return response, nil
}

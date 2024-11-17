package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

type UseCase interface {
	SignUp(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error)
	SignIn(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error)
	SignInVerify(ctx context.Context, in *authpb.SignInVerifyRequest) (*authpb.SignInVerifyResponse, error)
	GetUsersWithTokens(ctx context.Context, in *authpb.GetUsersWithTokensRequest) (*authpb.GetUsersWithTokensResponse, error)
	GetUsersWithProfiles(ctx context.Context, in *authpb.GetUsersWithProfilesRequest) (*authpb.GetUsersWithProfilesResponse, error)
	AddUserRole(ctx context.Context, in *authpb.AddUserRoleRequest) (*authpb.AddUserRoleResponse, error)
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

func (s *Server) GetUserByUUID(ctx context.Context, in *authpb.GetUserByUUIDRequest) (*authpb.GetUserByUUIDResponse, error) {
	panic("qwe")
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

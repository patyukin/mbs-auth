package server

import (
	"context"

	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

type UseCase interface {
	SignUpV1UseCase(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error)
	SignInV1UseCase(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error)
	SignInConfirmationV1UseCase(ctx context.Context, in *authpb.SignInConfirmationRequest) (*authpb.SignInConfirmationResponse, error)
	GetUserByIDUseCase(ctx context.Context, in *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error)
	GetBriefUserByID(ctx context.Context, in *authpb.GetBriefUserByIDRequest) (*authpb.GetBriefUserByIDResponse, error)
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

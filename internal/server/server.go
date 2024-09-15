package server

import (
	"context"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
)

type UseCase interface {
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

func (uc *Server) SignUp(context.Context, *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	panic("qwe")
}

func (uc *Server) SignIn(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	panic("qwe")
}

func (uc *Server) GetUserByUUID(ctx context.Context, in *authpb.GetUserByUUIDRequest) (*authpb.GetUserByUUIDResponse, error) {
	panic("qwe")
}

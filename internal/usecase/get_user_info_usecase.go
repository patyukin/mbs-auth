package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) GetUserInfoUseCase(ctx context.Context, in *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error) {
	userInfo, err := u.registry.GetRepo().SelectUserInfoByID(ctx, in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.GetRepo().SelectUserInfoByID: %w", err)
	}

	log.Debug().Msgf("userInfo: %v", userInfo)

	if userInfo == nil {
		return nil, fmt.Errorf("failed u.registry.GetRepo().SelectUserInfoByID: %w", errs.ErrUserNotFound)
	}

	return &authpb.GetUserByIDResponse{User: userInfo}, nil
}

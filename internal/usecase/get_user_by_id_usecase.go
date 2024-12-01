package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (u *UseCase) GetUserByIDUseCase(ctx context.Context, in *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error) {
	userInfo, err := u.registry.GetRepo().SelectUserInfoByID(ctx, in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.GetRepo().SelectUserInfoByID: %w", err)
	}

	pbm := model.ToProtoUserInfoDB(userInfo)

	return &authpb.GetUserByIDResponse{User: pbm}, nil
}

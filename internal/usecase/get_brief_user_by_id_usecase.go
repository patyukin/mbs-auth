package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (u *UseCase) GetBriefUserById(ctx context.Context, in *authpb.GetBriefUserByIDRequest) (*authpb.GetBriefUserByIDResponse, error) {
	user, err := u.registry.GetRepo().SelectBriefUserByUUID(ctx, in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.GetRepo().SelectUserByUUID: %w", err)
	}

	return model.ToProtoBriefUser(user), nil
}

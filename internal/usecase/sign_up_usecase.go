package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
)

func (u *UseCase) SignUp(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	var err error
	var uid uuid.UUID
	var user model.User
	var profile model.Profile
	var tgUser model.TelegramUser

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo db.RepositoryInterface) error {
		in.Password, err = u.HashPassword(in.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user = model.UserModelFromSignUpRequest(in)
		uid, err = repo.InsertIntoUsers(ctx, user)
		if err != nil {
			return fmt.Errorf("failed repo.InsertIntoUsers: %w", err)
		}

		profile, err = model.ProfileModelFromSignUpRequest(uid, in)
		if err != nil {
			return fmt.Errorf("failed model.ProfileModelFromSignUpRequest: %w", err)
		}

		_, err = repo.InsertIntoProfiles(ctx, profile)
		if err != nil {
			return fmt.Errorf("failed repo.InsertIntoProfiles: %w", err)
		}

		tgUser, err = model.TelegramUserModelFromSignUpRequest(uid, in)
		_, err = repo.InsertIntoTelegramUsers(ctx, tgUser)
		if err != nil {
			return fmt.Errorf("failed repo.InsertIntoTelegramUsers: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &authpb.SignUpResponse{UserId: uid.String()}, nil
}

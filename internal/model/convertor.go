package model

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"time"
)

func ProfileModelFromSignUpRequest(userUUID uuid.UUID, in *authpb.SignUpRequest) (Profile, error) {
	var Patronymic sql.NullString
	if in.Patronymic != "" {
		Patronymic.String = in.Patronymic
		Patronymic.Valid = true
	}

	layout := "02-01-2006"
	dateOfBirth, err := time.Parse(layout, in.DateOfBirth)
	if err != nil {
		return Profile{}, fmt.Errorf("failed time.Parse with in.DateOfBirth: %w", err)
	}

	return Profile{
		UserUUID:    userUUID,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Patronymic:  Patronymic,
		DateOfBirth: dateOfBirth,
		Email:       in.Email,
		Phone:       in.Phone,
		Address:     in.Address,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func UsersWithProfilesResponseFromUserWithProfile(u []UserWithProfile) []*authpb.UserGUWP {
	var users []*authpb.UserGUWP
	for _, v := range u {
		profile := &authpb.ProfileGUWP{
			FirstName:   v.FirstName,
			LastName:    v.LastName,
			Patronymic:  v.Patronymic.String,
			DateOfBirth: v.DateOfBirth.Format(time.DateOnly),
			Email:       v.ProfileEmail,
			Phone:       v.Phone,
			Address:     v.Address,
		}

		user := &authpb.UserGUWP{
			Id:      v.ID,
			Email:   v.Email,
			Role:    v.Role,
			Profile: profile,
		}

		users = append(users, user)
	}

	return users
}

func TelegramUserModelFromSignUpRequest(userUUID uuid.UUID, in *authpb.SignUpRequest) (TelegramUser, error) {
	return TelegramUser{
		UserUUID:      userUUID,
		TelegramLogin: in.TelegramLogin,
		CreatedAt:     time.Now().UTC(),
	}, nil
}

func UserModelFromSignUpRequest(in *authpb.SignUpRequest) User {
	return User{
		Email:        in.Email,
		PasswordHash: in.Password,
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	}
}

func SignUpRequestToUserModel(in *authpb.SignUpRequest) *User {
	return &User{
		UUID:         uuid.New(),
		Email:        in.Email,
		PasswordHash: in.Password,
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	}
}

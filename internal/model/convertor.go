package model

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"time"
)

func ProfileModelFromSignUpRequest(userUUID uuid.UUID, in *authpb.SignUpRequest) (Profile, error) {
	var Patronymic sql.NullString
	if in.Patronymic != "" {
		Patronymic.String = in.Patronymic
		Patronymic.Valid = true
	}

	layout := "2006-01-02"
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

func UserModelFromSignUpRequest(in *authpb.SignUpRequest) User {
	return User{
		Email:        in.Email,
		PasswordHash: in.Password,
		CreatedAt:    time.Now().UTC(),
	}
}

func ToProtoUserInfo(users []UserWithProfile) []*authpb.UserInfo {
	result := make([]*authpb.UserInfo, 0, len(users))
	for _, u := range users {
		result = append(
			result, &authpb.UserInfo{
				Id:    u.ID,
				Email: u.Email,
				Profile: &authpb.Profile{
					FirstName:   u.FirstName,
					LastName:    u.LastName,
					Patronymic:  u.Patronymic.String,
					DateOfBirth: u.DateOfBirth.Format("2006-01-02"),
					Phone:       u.Phone,
					Address:     u.Address,
				},
			},
		)
	}

	return result
}

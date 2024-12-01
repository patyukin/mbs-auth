package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func ProfileModelFromSignUpRequest(userUUID uuid.UUID, in *authpb.SignUpRequest) (Profile, error) {
	var Patronymic sql.NullString
	if in.GetPatronymic() != "" {
		Patronymic.String = in.GetPatronymic()
		Patronymic.Valid = true
	}

	layout := "2006-01-02"
	dateOfBirth, err := time.Parse(layout, in.GetDateOfBirth())
	if err != nil {
		return Profile{}, fmt.Errorf("failed time.Parse with in.DateOfBirth: %w", err)
	}

	return Profile{
		UserUUID:    userUUID,
		FirstName:   in.GetFirstName(),
		LastName:    in.GetLastName(),
		Patronymic:  Patronymic,
		DateOfBirth: dateOfBirth,
		Email:       in.GetEmail(),
		Phone:       in.GetPhone(),
		Address:     in.GetAddress(),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func UserModelFromSignUpRequest(in *authpb.SignUpRequest) User {
	return User{
		Email:        in.GetEmail(),
		PasswordHash: in.GetPassword(),
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

func ToProtoUserInfoDB(userInfoDB UserInfoDB) *authpb.UserInfo {
	return &authpb.UserInfo{
		Id:    userInfoDB.ID,
		Email: userInfoDB.Email,
		Profile: &authpb.Profile{
			FirstName:   userInfoDB.FirstName,
			LastName:    userInfoDB.LastName,
			Patronymic:  userInfoDB.Patronymic,
			DateOfBirth: userInfoDB.DateOfBirth,
			Phone:       userInfoDB.Phone,
			Address:     userInfoDB.Address,
		},
	}
}

func ToProtoBriefUser(user BriefUser) *authpb.GetBriefUserByIDResponse {
	return &authpb.GetBriefUserByIDResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ChatId:    user.ChatID,
	}
}

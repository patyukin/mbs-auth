package model

import (
	"github.com/google/uuid"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"time"
)

func SignUpRequestToUserModel(in *authpb.SignUpRequest) *User {
	return &User{
		UUID:         uuid.New(),
		Email:        in.Email,
		PasswordHash: in.Password,
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	}
}

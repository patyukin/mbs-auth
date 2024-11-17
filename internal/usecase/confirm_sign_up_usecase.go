package usecase

import (
	"context"
	"github.com/patyukin/mbs-pkg/pkg/model"
)

func (u *UseCase) ConfirmSignUp(ctx context.Context, msg model.AuthSignUpConfirmCode) error {
	return nil
}

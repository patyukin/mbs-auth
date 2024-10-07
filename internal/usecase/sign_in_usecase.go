package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"time"
)

func (u *UseCase) SignIn(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	var err error
	var user model.User
	var telegramUser model.TelegramUser
	var msg []byte

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo db.RepositoryInterface) error {
		user, err = repo.SelectUserByEmail(ctx, in.Email)
		if err != nil {
			return fmt.Errorf("failed to select user: %w", err)
		}

		err = u.ComparePasswords([]byte(user.PasswordHash), in.Password)
		if err != nil {
			return fmt.Errorf("failed to compare passwords: %w", err)
		}

		// Генерация уникального кода 2FA
		var code string
		var exists int64
		for {
			code, err = u.GenerateSignInCode()
			if err != nil {
				return fmt.Errorf("failed to generate sign in code: %w", err)
			}

			// Проверка на уникальность кода
			exists, err = u.chr.Exists2FACode(ctx, code)
			if err != nil {
				return fmt.Errorf("failed to check sign in code: %w", err)
			}

			if exists == 0 {
				break
			}
		}

		err = u.chr.Set2FACode(ctx, code, user.UUID.String(), 5*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to set 2fa code: %w", err)
		}

		telegramUser, err = repo.SelectFromTelegramUsersByUser(ctx, user.UUID)
		if err != nil {
			return fmt.Errorf("failed repo.SelectFromTelegramUsersByUser: %w", err)
		}

		if !telegramUser.TelegramChatID.Valid {
			return fmt.Errorf("telegram chat id not found")
		}

		payload := struct {
			ChatId int64  `json:"chat_id"`
			Code   string `json:"code"`
		}{
			Code:   code,
			ChatId: telegramUser.TelegramChatID.Int64,
		}

		msg, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		err = u.prdcr.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("failed u.prdcr.SendMessage: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read committed: %w", err)
	}

	return nil, nil
}

package cacher

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Cacher struct {
	client *redis.Client
}

func New(ctx context.Context, dsn string) (*Cacher, error) {
	cl := redis.NewClient(&redis.Options{Addr: dsn})

	err := cl.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Msg("connected to redis")
	return &Cacher{client: cl}, nil
}

func (r *Cacher) SetVerificationCode(ctx context.Context, telegramUserID int, code string, expiration time.Duration) error {
	err := r.client.Set(ctx, fmt.Sprintf("code:%d", telegramUserID), code, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set verification code: %w", err)
	}

	return nil
}

func (r *Cacher) GetVerificationCode(ctx context.Context, telegramUserID int) (string, error) {
	code, err := r.client.Get(ctx, fmt.Sprintf("code:%d", telegramUserID)).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}

	return code, nil
}

func (r *Cacher) SetSignUpCode(ctx context.Context, tgUserName string, code, userUUID uuid.UUID, expiration time.Duration) error {
	err := r.client.Set(ctx, "user:"+tgUserName, fmt.Sprintf("%s:%s", code.String(), userUUID.String()), expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set sign up code: %w", err)
	}

	return nil
}

func (r *Cacher) GetSignUpCode(ctx context.Context, tgUserName string) (string, error) {
	code, err := r.client.Get(ctx, "user:"+tgUserName).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get sign up code: %w", err)
	}

	return code, nil
}

func (r *Cacher) DeleteSignUpCode(ctx context.Context, tgUserName string) error {
	err := r.client.Del(ctx, "user:"+tgUserName).Err()
	if err != nil {
		return fmt.Errorf("failed to delete sign up code: %w", err)
	}

	return nil
}

func (r *Cacher) Set2FACode(ctx context.Context, userID, code string) error {
	err := r.client.Set(ctx, fmt.Sprintf("otp2fa:%s", userID), code, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set 2fa code: %w", err)
	}

	return nil
}

func (r *Cacher) Get2FACode(ctx context.Context, userID string) (string, error) {
	code, err := r.client.Get(ctx, fmt.Sprintf("otp2fa:%s", userID)).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get 2fa code: %w", err)
	}

	return code, nil
}

func (r *Cacher) Delete2FACode(ctx context.Context, userID string) error {
	err := r.client.Del(ctx, fmt.Sprintf("otp2fa:%s", userID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete 2fa code: %w", err)
	}

	return nil
}

func (r *Cacher) SetTempCode(ctx context.Context, userID, code string, expiration time.Duration) error {
	err := r.client.Set(ctx, "tempcode:"+userID, code, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set temp code: %w", err)
	}

	return nil
}

func (r *Cacher) GetTempCode(ctx context.Context, userID string) (string, error) {
	code, err := r.client.Get(ctx, "tempcode:"+userID).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get temp code: %w", err)
	}

	return code, nil
}

func (r *Cacher) DeleteTempCode(ctx context.Context, userID string) error {
	err := r.client.Del(ctx, "tempcode:"+userID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete temp code: %w", err)
	}

	return nil
}

func (r *Cacher) Close() error {
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close redis: %w", err)
	}

	return nil
}

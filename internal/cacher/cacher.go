package cacher

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
	"time"
)

type Cacher struct {
	client *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (*Cacher, error) {
	c := redis.NewClient(&redis.Options{Addr: net.JoinHostPort(cfg.Redis.Host, strconv.Itoa(cfg.Redis.Port))})

	err := c.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Msg("connected to redis")
	return &Cacher{client: c}, nil
}

func (r *Cacher) SetVerificationCode(ctx context.Context, telegramUserID int, code string, expiration time.Duration) error {
	return r.client.Set(ctx, fmt.Sprintf("code:%d", telegramUserID), code, expiration).Err()
}

func (r *Cacher) GetVerificationCode(ctx context.Context, telegramUserID int) (string, error) {
	return r.client.Get(ctx, fmt.Sprintf("code:%d", telegramUserID)).Result()
}

func (r *Cacher) SetSignUpCode(ctx context.Context, tgUserName string, code, userUUID uuid.UUID, expiration time.Duration) error {
	return r.client.Set(ctx, "user:"+tgUserName, fmt.Sprintf("%s:%s", code.String(), userUUID.String()), expiration).Err()
}

func (r *Cacher) GetSignUpCode(ctx context.Context, tgUserName string) (string, error) {
	return r.client.Get(ctx, "user:"+tgUserName).Result()
}

func (r *Cacher) DeleteSignUpCode(ctx context.Context, tgUserName string) error {
	return r.client.Del(ctx, "user:"+tgUserName).Err()
}

func (r *Cacher) Exists2FACode(ctx context.Context, code string) (int64, error) {
	return r.client.Exists(ctx, code).Result()
}

func (r *Cacher) Set2FACode(ctx context.Context, code, userID string, expiration time.Duration) error {
	return r.client.Set(ctx, code, userID, expiration).Err()
}

func (r *Cacher) Get2FACode(ctx context.Context, code string) (string, error) {
	return r.client.Get(ctx, code).Result()
}

func (r *Cacher) Delete2FACode(ctx context.Context, code string) error {
	return r.client.Del(ctx, code).Err()
}

func (r *Cacher) SetTempCode(ctx context.Context, userID string, code string, expiration time.Duration) error {
	return r.client.Set(ctx, "tempcode:"+userID, code, expiration).Err()
}

func (r *Cacher) GetTempCode(ctx context.Context, userID string) (string, error) {
	return r.client.Get(ctx, "tempcode:"+userID).Result()
}

func (r *Cacher) DeleteTempCode(ctx context.Context, userID string) error {
	return r.client.Del(ctx, "tempcode:"+userID).Err()
}

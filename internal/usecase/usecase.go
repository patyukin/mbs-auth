package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-auth/internal/telegram"
	"time"
)

type Registry interface {
	GetRepo() db.RepositoryInterface
	ReadCommitted(ctx context.Context, f db.Handler) error
}

type Producer interface {
	SendMessage(body []byte) error
}

type Cacher interface {
	SetSignUpCode(ctx context.Context, tgUserName string, code, userUUID uuid.UUID, expiration time.Duration) error
	Exists2FACode(ctx context.Context, code string) (int64, error)
	Set2FACode(ctx context.Context, code, userID string, expiration time.Duration) error
	Get2FACode(ctx context.Context, code string) (string, error)
}

type UseCase struct {
	registry  Registry
	prdcr     Producer
	chr       Cacher
	bot       *telegram.Bot
	jwtSecret []byte
}

func New(registry Registry, prdcr Producer, chr Cacher, bot *telegram.Bot, jwtSecret string) *UseCase {
	return &UseCase{
		registry:  registry,
		prdcr:     prdcr,
		chr:       chr,
		bot:       bot,
		jwtSecret: []byte(jwtSecret),
	}
}

func (u *UseCase) GetJWTToken() []byte {
	return u.jwtSecret
}

func (u *UseCase) GetTelegramBot() string {
	return u.bot.API.Self.UserName
}
func (u *UseCase) GenerateSignInCode() (string, error) {
	bytes := make([]byte, 30)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

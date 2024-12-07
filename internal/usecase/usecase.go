package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/config"
	"github.com/patyukin/mbs-auth/internal/db"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrCodesNotEqual          = errors.New("codes are not equal")
	ErrTelegramChatIDNotFound = errors.New("telegram chat id not found")
)

const maxBytes = 30

type Producer interface {
	EnqueueTelegramMessage(ctx context.Context, body []byte, headers amqp.Table) error
}

type Cacher interface {
	SetSignUpCode(ctx context.Context, tgUserName string, code, userUUID uuid.UUID, expiration time.Duration) error
	GetSignUpCode(ctx context.Context, tgUserName string) (string, error)
	DeleteSignUpCode(ctx context.Context, tgUserName string) error
	Set2FACode(ctx context.Context, userID, code string) error
	Get2FACode(ctx context.Context, userID string) (string, error)
	Delete2FACode(ctx context.Context, userID string) error
}

type UseCase struct {
	registry  *db.Registry
	prdcr     Producer
	chr       Cacher
	bot       string
	jwtSecret []byte
}

func New(registry *db.Registry, prdcr Producer, chr Cacher, cfg *config.Config) *UseCase {
	return &UseCase{
		registry:  registry,
		prdcr:     prdcr,
		chr:       chr,
		bot:       cfg.TelegramBotName,
		jwtSecret: []byte(cfg.JwtSecret),
	}
}

func (u *UseCase) GetJWTToken() []byte {
	return u.jwtSecret
}

func (u *UseCase) GetTelegramBot() string {
	return u.bot
}

func (u *UseCase) GenerateSignInCode() (string, error) {
	bytes := make([]byte, maxBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func (u *UseCase) generateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedString, nil
}

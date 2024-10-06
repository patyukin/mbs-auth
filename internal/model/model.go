package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID         uuid.UUID    `json:"id"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
	Role         string       `json:"role"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
}

type TelegramUser struct {
	UUID           uuid.UUID     `json:"id"`
	UserUUID       uuid.UUID     `json:"user_id"`
	TelegramLogin  string        `json:"telegram_login"`
	TelegramUserID sql.NullInt64 `json:"telegram_id"`
	TelegramChatID sql.NullInt64 `json:"chat_id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      sql.NullTime  `json:"updated_at"`
}

type Token struct {
	Token     uuid.UUID    `json:"token"`
	UserUUID  uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type Profile struct {
	UUID        uuid.UUID      `json:"id"`
	UserUUID    uuid.UUID      `json:"user_id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Patronymic  sql.NullString `json:"patronymic"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   sql.NullTime   `json:"updated_at"`
}

type UserWithProfile struct {
	ID           string
	Email        string
	Role         string
	FirstName    string
	LastName     string
	Patronymic   sql.NullString
	DateOfBirth  time.Time
	ProfileEmail string
	Phone        string
	Address      string
}

package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/model"
)

func (r *Repository) InsertIntoTelegramUsers(ctx context.Context, in model.TelegramUser) (uuid.UUID, error) {
	query := `INSERT INTO telegram_users (user_id, telegram_login, created_at) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.UserUUID, in.TelegramLogin, in.CreatedAt)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed r.db.QueryRowContext, row.Err(): %w", row.Err())
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return id, nil
}

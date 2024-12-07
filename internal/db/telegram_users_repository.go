package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/model"
)

func (r *Repository) InsertIntoTelegramUsers(ctx context.Context, userUUID uuid.UUID, telegramLogin string) (uuid.UUID, error) {
	query := `INSERT INTO telegram_users (user_id, telegram_login, created_at) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, userUUID.String(), telegramLogin, time.Now().UTC())
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

func (r *Repository) SelectFromTelegramUsersByUser(ctx context.Context, userUUID uuid.UUID) (model.TelegramUser, error) {
	query := `SELECT id, user_id,  telegram_id, chat_id, created_at, updated_at FROM telegram_users WHERE user_id = $1`
	row := r.db.QueryRowContext(ctx, query, userUUID.String())
	if row.Err() != nil {
		return model.TelegramUser{}, fmt.Errorf("failed to select telegram_users: %w", row.Err())
	}

	var tgUser model.TelegramUser
	err := row.Scan(&tgUser.UUID, &tgUser.UserUUID, &tgUser.TelegramUserID, &tgUser.TelegramChatID, &tgUser.CreatedAt, &tgUser.UpdatedAt)
	if err != nil {
		return model.TelegramUser{}, fmt.Errorf("failed to select telegram_users: %w", err)
	}

	return tgUser, nil
}

func (r *Repository) RemoveTelegramUsersByUserIDs(ctx context.Context, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]any, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM telegram_users WHERE user_id IN (%s)", strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed removing users by IDs in telegram_users: %w", err)
	}

	return nil
}

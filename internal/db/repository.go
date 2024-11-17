package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Repository struct {
	db QueryExecutor
}

func (r *Repository) UpdateTelegramUserAfterSignUp(ctx context.Context, userUUID uuid.UUID, chatID, userID int64) error {
	currentTime := time.Now().UTC()
	query := `UPDATE telegram_users SET telegram_id = $1, chat_id = $2, updated_at = $3 WHERE user_id = $4`
	_, err := r.db.ExecContext(ctx, query, userID, chatID, currentTime, userUUID.String())
	if err != nil {
		return fmt.Errorf("failed to insert in telegram_users: %w", err)
	}

	return nil
}

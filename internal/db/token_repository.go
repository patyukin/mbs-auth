package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *Repository) InsertToken(ctx context.Context, userUUID uuid.UUID) (string, error) {
	currentTime := time.Now().UTC()
	expiresAt := currentTime.Add(24 * 30 * time.Hour)
	query := `INSERT INTO tokens (user_id, expires_at, created_at) VALUES ($1, $2, $3) RETURNING token`
	row := r.db.QueryRowContext(ctx, query, userUUID.String(), expiresAt, currentTime)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to insert token: %w", row.Err())
	}

	var token uuid.UUID
	err := row.Scan(&token)
	if err != nil {
		return "", fmt.Errorf("failed to insert token: %w", err)
	}

	return token.String(), nil
}

func (r *Repository) CleanExpiredTokens(ctx context.Context) error {
	currentTime := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, "DELETE FROM tokens WHERE expires_at < $1", currentTime)
	if err != nil {
		return fmt.Errorf("failed cleaning tokens: %w", err)
	}

	return nil
}

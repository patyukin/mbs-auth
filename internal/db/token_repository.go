package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *Repository) UpsertToken(ctx context.Context, userUUID uuid.UUID) (string, error) {
	currentTime := time.Now().UTC()
	expiresAt := currentTime.Add(24 * 30 * time.Hour)
	query := `
INSERT INTO tokens (user_id, expires_at, created_at, token)
VALUES ($1, $2, $3, DEFAULT)
ON CONFLICT (user_id) 
DO UPDATE SET
    expires_at = EXCLUDED.expires_at,
    created_at = EXCLUDED.created_at,
    token = DEFAULT
RETURNING token
`
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

func (r *Repository) SelectByID(ctx context.Context, ID string) (string, error) {
	timeNow := time.Now().UTC()
	query := `SELECT user_id FROM tokens WHERE token = $1 AND expires_at > $2`
	row := r.db.QueryRowContext(ctx, query, ID, timeNow)
	if row.Err() != nil {
		return "", fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed row.Scan: %w", err)
	}

	return userID, nil
}

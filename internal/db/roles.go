package db

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) InsertIntoUsersRoles(ctx context.Context, userID, roleID string) error {
	currentTime := time.Now().UTC()
	query := `INSERT INTO users_roles (user_id, role_id, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, userID, roleID, currentTime)
	if err != nil {
		return fmt.Errorf("failed to insert into users_roles: %w", err)
	}

	return nil
}

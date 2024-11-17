package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/model"
	"strings"
)

func (r *Repository) InsertIntoProfiles(ctx context.Context, in model.Profile) (uuid.UUID, error) {
	query := `
INSERT INTO profiles (user_id, first_name, last_name, patronymic, date_of_birth, email, phone, address, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := r.db.QueryRowContext(
		ctx,
		query,
		in.UserUUID,
		in.FirstName,
		in.LastName,
		in.Patronymic,
		in.DateOfBirth,
		in.Email,
		in.Phone,
		in.Address,
		in.CreatedAt,
	)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return id, nil
}

func (r *Repository) RemoveProfilesByUserIDs(ctx context.Context, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]any, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM profiles WHERE user_id IN (%s)", strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed removing users by IDs in telegram_users: %w", err)
	}

	return nil
}

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/model"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (r *Repository) InsertIntoUsers(ctx context.Context, in model.User) (uuid.UUID, error) {
	query := `INSERT INTO users (email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.Email, in.PasswordHash, in.CreatedAt)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return id, nil
}

func (r *Repository) SelectUsersWithTokensCount(ctx context.Context) (int32, error) {
	query := `SELECT COUNT(*) FROM users u INNER JOIN tokens t ON u.id = t.user_id`
	row := r.db.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var count int32
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed row.Scan: %w", err)
	}

	return count, nil
}

func (r *Repository) SelectUsersWithTokens(ctx context.Context, limit int32, page int32) ([]*authpb.UserGUWR, error) {
	query := `
SELECT
    u.id,
    u.email,
    json_agg(
        json_build_object(
            'token', t.token,
            'expires_at', t.expires_at
        )
    ) AS tokens
FROM (
    SELECT
        u_inner.id,
        u_inner.email,
        ROW_NUMBER() OVER (ORDER BY u_inner.created_at ASC) as rn
    FROM users u_inner
    INNER JOIN tokens t_inner ON u_inner.id = t_inner.user_id
) u
INNER JOIN tokens t ON u.id = t.user_id
WHERE u.rn > ($1 - 1) * $2 AND u.rn <= $1 * $2
GROUP BY u.id, u.email, u.rn
ORDER BY u.rn ASC
`
	rows, err := r.db.QueryContext(ctx, query, limit, page)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error().Msgf("failed rows.Close: %v", err)
		}
	}(rows)

	var users []*authpb.UserGUWR
	for rows.Next() {
		var id, email string
		var tokensJSON []byte

		if err = rows.Scan(&id, &email, &tokensJSON); err != nil {
			return nil, fmt.Errorf("rows.Scan tokensJSON: %w", err)
		}

		var tokens []*authpb.TokenGUWR

		if err = json.Unmarshal(tokensJSON, &tokens); err != nil {
			return nil, fmt.Errorf("json.Unmarshal tokens: %w", err)
		}

		user := &authpb.UserGUWR{
			Id:     id,
			Email:  email,
			Tokens: tokens,
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) SelectUsersWithProfilesCount(ctx context.Context) (int32, error) {
	query := `SELECT COUNT(*) FROM users u INNER JOIN profiles p ON u.id = p.user_id`
	row := r.db.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var count int32
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed row.Scan: %w", err)
	}

	return count, nil
}

func (r *Repository) SelectUsersWithProfiles(ctx context.Context, limit int32, page int32) ([]model.UserWithProfile, error) {
	offset := (int(page) - 1) * int(limit)

	query := `
SELECT
	u.id,
	u.email,
	p.first_name,
	p.last_name,
	p.patronymic,
	p.date_of_birth,
	p.email,
	p.phone,
	p.address
FROM users u
INNER JOIN profiles p ON u.id = p.user_id
ORDER BY u.created_at ASC
OFFSET $1 LIMIT $2;
`
	rows, err := r.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed rows.Err(): %w", err)
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error().Msgf("failed rows.Close: %v", err)
		}
	}(rows)

	var uwps []model.UserWithProfile

	for rows.Next() {
		var uwp model.UserWithProfile
		err = rows.Scan(
			&uwp.ID,
			&uwp.Email,
			&uwp.FirstName,
			&uwp.LastName,
			&uwp.Patronymic,
			&uwp.DateOfBirth,
			&uwp.ProfileEmail,
			&uwp.Phone,
			&uwp.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("failed TempUser rows.Scan: %w", err)
		}

		uwps = append(uwps, uwp)
	}

	return uwps, nil
}

func (r *Repository) SelectRegisteredUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
SELECT
    u.id,
    u.email,
    u.password_hash
FROM users AS u
INNER JOIN telegram_users AS tu ON u.id = tu.user_id 
WHERE email = $1 AND tu.chat_id IS NOT NULL`
	row := r.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.UUID, &user.Email, &user.PasswordHash)
	if err != nil {
		return model.User{}, fmt.Errorf("failed row.Scan: %w", errs.ErrUserNotFound)
	}

	return user, nil
}

func (r *Repository) SelectUserByUUID(ctx context.Context, userUUID string) (model.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userUUID)
	if row.Err() != nil {
		return model.User{}, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var user model.User
	err := row.Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return user, nil
}

func (r *Repository) SelectNotRegisteredUsers(ctx context.Context) ([]uuid.UUID, error) {
	// TODO uncomment it
	t := time.Now().UTC() //.Add(-2 * time.Hour)
	query := `
SELECT 
  u.id 
FROM users AS u 
	INNER JOIN telegram_users AS tu ON u.id = tu.user_id 
WHERE tu.chat_id IS NULL 
  AND tu.created_at < $1
`

	rows, err := r.db.QueryContext(ctx, query, t)
	if err != nil {
		return nil, fmt.Errorf("failed to select users in r.db.QueryContext: %w", err)
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error().Msgf("failed rows.Close: %v", err)
		}
	}(rows)

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration in rows.Err(): %w", err)
	}

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user ID in rows.Scan: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repository) RemoveUsersByIDs(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM users WHERE id IN (%s)", strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed removing users by IDs: %w", err)
	}

	return nil
}

func (r *Repository) SelectUserWithExistsEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT id FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	if row.Err() != nil {
		return false, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed row.Scan: %w", err)
	}

	return true, nil
}

func (r *Repository) SelectUserInfoByID(ctx context.Context, userID string) (*authpb.UserInfo, error) {
	query := `
SELECT
		u.id,
		u.email,
		p.first_name,
		p.last_name,
		p.patronymic,
		p.date_of_birth,
		p.phone,
		p.address,
		tu.telegram_login,
		tu.telegram_id,
		tu.chat_id,
		TO_CHAR(u.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
FROM users AS u
INNER JOIN telegram_users AS tu ON u.id = tu.user_id
INNER JOIN profiles AS p ON u.id = p.user_id
WHERE u.id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var user authpb.UserInfo
	err := row.Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.DateOfBirth,
		&user.Phone,
		&user.Address,
		&user.TelegramLogin,
		&user.TelegramId,
		&user.ChatId,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed row.Scan: %w", err)
	}

	return &user, nil
}

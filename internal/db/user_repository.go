package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/model"
	authpb "github.com/patyukin/mbs-auth/pkg/auth_v1"
	"github.com/rs/zerolog/log"
)

func (r *Repository) InsertIntoUsers(ctx context.Context, in model.User) (uuid.UUID, error) {
	query := `INSERT INTO users (email, password_hash, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.Email, in.PasswordHash, in.Role, in.CreatedAt)
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
    u.role,
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
        u_inner.role,
        ROW_NUMBER() OVER (ORDER BY u_inner.created_at ASC) as rn
    FROM users u_inner
    INNER JOIN tokens t_inner ON u_inner.id = t_inner.user_id
) u
INNER JOIN tokens t ON u.id = t.user_id
WHERE u.rn > ($1 - 1) * $2 AND u.rn <= $1 * $2
GROUP BY u.id, u.email, u.role, u.rn
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
		var id, email, role string
		var tokensJSON []byte

		if err = rows.Scan(&id, &email, &role, &tokensJSON); err != nil {
			return nil, fmt.Errorf("rows.Scan tokensJSON: %w", err)
		}

		var tokens []*authpb.TokenGUWR

		if err = json.Unmarshal(tokensJSON, &tokens); err != nil {
			return nil, fmt.Errorf("json.Unmarshal tokens: %w", err)
		}

		user := &authpb.UserGUWR{
			Id:     id,
			Email:  email,
			Role:   role,
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
	u.role,
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
			&uwp.Role,
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

func (r *Repository) SelectUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return model.User{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return user, nil
}

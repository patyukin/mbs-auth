package db

import (
	"context"
	"fmt"

	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

func (r *Repository) InsertIntoUsersRoles(ctx context.Context, userID, roleID string) error {
	query := `INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to insert into users_roles: %w", err)
	}

	return nil
}

func (r *Repository) SelectExistsRowByRoleUserIDAndRoutePath(ctx context.Context, in *authpb.AuthorizeUserRequest) (bool, error) {
	log.Debug().Msgf("in: %v", in)
	query := `
SELECT
  ur.id
FROM users_roles AS ur
INNER JOIN roles_permissions AS rp ON ur.role_id = rp.role_id
INNER JOIN roles AS r ON rp.role_id = r.id
INNER JOIN permissions AS p ON rp.permission_id = p.id
WHERE ur.user_id = $1
	AND $2 ~ ('^' || regexp_replace(p.route_path, '\{[^}]+\}', '[^/]+', 'g') || '$')
  AND p.method = $3`

	rows := r.db.QueryRowContext(ctx, query, in.GetUserId(), in.GetRoutePath(), in.GetMethod())
	if rows.Err() != nil {
		return false, fmt.Errorf("failed r.db.QueryRowContext: %w", rows.Err())
	}

	var id string
	err := rows.Scan(&id)
	if err != nil {
		return false, fmt.Errorf("failed rows.Scan: %w", err)
	}

	if id == "" {
		return false, fmt.Errorf("failed ID is empty: %w", err)
	}

	return true, nil
}

func (r *Repository) SelectRoleByUserID(ctx context.Context, userID string) (string, error) {
	query := `SELECT role_id FROM users_roles WHERE user_id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
	if row.Err() != nil {
		return "", fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var roleID string
	err := row.Scan(&roleID)
	if err != nil {
		return "", fmt.Errorf("failed row.Scan: %w", err)
	}

	return roleID, nil
}

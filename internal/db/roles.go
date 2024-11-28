package db

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
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

func (r *Repository) SelectExistsRowByRoleUserIDAndRoutePath(ctx context.Context, in *authpb.AuthorizeUserRequest) (bool, error) {
	query := `
SELECT
	ur.id
FROM users_roles AS ur
INNER JOIN roles_permissions AS rp ON ur.role_id = rp.role_id
INNER JOIN roles AS r ON rp.role_id = r.id
INNER JOIN permissions AS p ON rp.permission_id = p.id
WHERE 
	ur.user_id = $1 
  AND p.route_path = $2 
  AND p.method = $3`

	rows := r.db.QueryRowContext(ctx, query, in.UserId, in.RoutePath, in.Method)
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

package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (r *userRepository) ListUsers(ctx context.Context, params model.ListUsersParams) ([]*domain.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, avatar_key, role, manager_id, created_at, updated_at FROM users`
	var (
		args       []any
		conditions []string
		argPos     = 1
	)

	if params.Role != nil {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argPos))
		args = append(args, *params.Role)
		argPos++
	}
	if params.Search != nil && *params.Search != "" {
		conditions = append(conditions, fmt.Sprintf(`(email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)`, argPos, argPos, argPos))
		args = append(args, "%"+*params.Search+"%")
		argPos++
	}
	if len(conditions) != 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(`ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argPos, argPos+1)
	args = append(args, params.Limit, params.Offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository.ListUsers: %w", err)
	}
	defer rows.Close()
	users := make([]*domain.User, 0)
	for rows.Next() {
		var m model.UserModel
		err = rows.Scan(
			&m.ID,
			&m.Email,
			&m.PasswordHash,
			&m.FirstName,
			&m.LastName,
			&m.AvatarKey,
			&m.Role,
			&m.ManagerID,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository.ListUsers scan: %w", err)
		}
		users = append(users, converter.ToUserEntity(&m))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.ListUsers rows: %w", err)
	}

	return users, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, userID uuid.UUID, role domain.Role) error {
	query := `
	UPDATE users
	SET role = $1,
		updated_at = NOW()
	WHERE id = $2`

	cmd, err := r.pool.Exec(ctx, query, role, userID)
	if err != nil {
		return fmt.Errorf("repository.UpdateRole(user): %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) CountAdmins(ctx context.Context) (int, error) {
	var count int

	query := `
	SELECT COUNT(*)
	FROM users
	WHERE role = 'admin'
	`

	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository.CountAdmins(user): %w", err)
	}

	return count, nil
}

func (r *userRepository) AssignManager(ctx context.Context, userID uuid.UUID, managerID uuid.UUID) error {
	query := `
	UPDATE users
	SET manager_id = $1,
		updated_at = NOW()
	WHERE id = $2`

	cmd, err := r.pool.Exec(ctx, query, managerID, userID)
	if err != nil {
		return fmt.Errorf("repository.AssignManager(user): %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) RemoveManager(ctx context.Context, userID uuid.UUID) error {
	query := `
	UPDATE users
	SET manager_id = NULL,
    	updated_at = NOW()
	WHERE id = $1
  		AND manager_id IS NOT NULL`

	cmd, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("repository.RemoveManager(user): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) ListEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.User, error) {

	query := `
	SELECT
		id,
		email,
		password_hash,
		first_name,
		last_name,
		avatar_key,
		role,
		manager_id,
		created_at,
		updated_at
	FROM users
	WHERE manager_id = $1
	ORDER BY last_name, first_name`

	rows, err := r.pool.Query(ctx, query, managerID)
	if err != nil {
		return nil, fmt.Errorf("repository.ListEmployeesByManager(user): %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0)

	for rows.Next() {
		var m model.UserModel
		err = rows.Scan(
			&m.ID,
			&m.Email,
			&m.PasswordHash,
			&m.FirstName,
			&m.LastName,
			&m.AvatarKey,
			&m.Role,
			&m.ManagerID,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository.ListEmployeesByManager(user): %w", err)
		}

		users = append(users, converter.ToUserEntity(&m))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.ListEmployeesByManager(user): %w", err)
	}

	return users, nil
}

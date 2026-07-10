package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	m := converter.ToUserModel(u)
	query := `INSERT INTO users(email, password_hash, first_name, last_name, role, manager_id) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	err := r.pool.QueryRow(ctx, query, m.Email, m.PasswordHash, m.FirstName, m.LastName, m.Role, m.ManagerID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repository.Create (user): %w", err)
	}
	return id, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m model.UserModel

	query := `
	SELECT id, email, password_hash, first_name, last_name, avatar_key, role, manager_id, created_at, updated_at
	FROM users
	WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("repository.GetByEmail(user): %w", err)
	}

	return converter.ToUserEntity(&m), nil
}

func (r *userRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m model.UserModel
	query := `SELECT id, email, password_hash, first_name, last_name, avatar_key, role, manager_id, created_at, updated_at
	FROM users
	WHERE id=$1`
	if err := r.pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.Email, &m.PasswordHash, &m.FirstName, &m.LastName, &m.AvatarKey, &m.Role, &m.ManagerID, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("repository.GetByID(user): %w", err)
	}
	return converter.ToUserEntity(&m), nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, update *domain.UpdateUserProfile) error {
	args := []any{}
	setClauses := []string{}
	i := 1

	if update.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email=$%d", i))
		args = append(args, *update.Email)
		i++
	}

	if update.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name=$%d", i))
		args = append(args, *update.FirstName)
		i++
	}

	if update.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name=$%d", i))
		args = append(args, *update.LastName)
		i++
	}

	setClauses = append(setClauses, "updated_at=NOW()")

	args = append(args, userID)

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id=$%d",
		strings.Join(setClauses, ", "),
		i,
	)

	cmd, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("repository.UpdateProfile(user): %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashPassword string) error {
	query := `UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`
	_, err := r.pool.Exec(ctx, query, hashPassword, userID)
	if err != nil {
		return fmt.Errorf("repository.UpdatePassword(user): %w", err)
	}
	return nil
}

func (r *userRepository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarKey *string) error {
	query := `
	UPDATE users
	SET avatar_key = $1,
	    updated_at = NOW()
	WHERE id = $2`

	cmd, err := r.pool.Exec(ctx, query, avatarKey, userID)
	if err != nil {
		return fmt.Errorf("repository.UpdateAvatar(user): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) ExistsAdmin(ctx context.Context) (bool, error) {
	var exists bool

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE role = 'admin'
	)
	`

	err := r.pool.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("repository.ExistsAdmin(user): %w", err)
	}

	return exists, nil
}

func (r *userRepository) ClearManagerAssignments(ctx context.Context, managerID uuid.UUID) error {
	query := `
		UPDATE users
		SET manager_id = NULL
		WHERE manager_id = $1
	`

	_, err := r.pool.Exec(ctx, query, managerID)
	return err
}

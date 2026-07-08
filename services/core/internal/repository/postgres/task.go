package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTaskNotFound = errors.New("task not found")

type taskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *taskRepository {
	return &taskRepository{
		pool: pool,
	}
}

func (r *taskRepository) GetNextPosition(ctx context.Context, planID uuid.UUID) (int, error) {
	var position int

	query := `
	SELECT COALESCE(MAX(position),0)+1
	FROM tasks
	WHERE plan_id=$1
	`

	err := r.pool.QueryRow(ctx, query, planID).Scan(&position)
	if err != nil {
		return 0, fmt.Errorf("repository.GetNextPosition(task): %w", err)
	}

	return position, nil
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	var id uuid.UUID

	m := converter.ToTaskModel(task)

	query := `
	INSERT INTO tasks(
		plan_id,
		title,
		description,
		position,
		status
	)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		m.PlanID,
		m.Title,
		m.Description,
		m.Position,
		m.Status,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("repository.Create(task): %w", err)
	}

	return id, nil
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var m model.TaskModel

	query := `
		SELECT
			id,
			plan_id,
			title,
			description,
			position,
			status,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.PlanID,
		&m.Title,
		&m.Description,
		&m.Position,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("repository.GetByID(task): %w", err)
	}

	return converter.ToTaskEntity(&m), nil
}

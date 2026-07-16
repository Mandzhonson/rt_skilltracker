package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTaskNotFound = errors.New("task not found")

type taskRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewTaskRepository(pool *pgxpool.Pool, log *slog.Logger) *taskRepository {
	return &taskRepository{
		pool: pool,
		log:  log,
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
		r.log.Error("Failed to get next position for task",
			slog.String("error", err.Error()),
			slog.String("plan_id", planID.String()),
		)
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
		r.log.Error("Failed to create task",
			slog.String("error", err.Error()),
			slog.String("plan_id", m.PlanID.String()),
			slog.String("title", m.Title),
		)
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

		r.log.Error("Failed to get task by ID",
			slog.String("error", err.Error()),
			slog.String("task_id", id.String()),
		)
		return nil, fmt.Errorf("repository.GetByID(task): %w", err)
	}

	return converter.ToTaskEntity(&m), nil
}

func (r *taskRepository) Update(ctx context.Context, id uuid.UUID, title *string, description *string) error {

	query := `
	UPDATE tasks
	SET
		title = COALESCE($1,title),
		description = $2,
		updated_at = NOW()
	WHERE id=$3
	`

	result, err := r.pool.Exec(ctx, query, title, description, id)
	if err != nil {
		r.log.Error("Failed to update task",
			slog.String("error", err.Error()),
			slog.String("task_id", id.String()),
		)
		return fmt.Errorf("repository.Update(task): %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {

	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.log.Error("Failed to delete task",
			slog.String("error", err.Error()),
			slog.String("task_id", id.String()),
		)
		return fmt.Errorf(
			"repository.Delete(task): %w",
			err,
		)
	}
	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {

	query := `
	UPDATE tasks
	SET
		status=$1,
		updated_at=NOW()
	WHERE id=$2
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		r.log.Error("Failed to update task status",
			slog.String("error", err.Error()),
			slog.String("task_id", id.String()),
			slog.String("status", status),
		)
		return fmt.Errorf("repository.UpdateStatus(task): %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *taskRepository) ListByPlanID(ctx context.Context, planID uuid.UUID) ([]*domain.Task, error) {

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
        WHERE plan_id = $1
        ORDER BY position
    `

	rows, err := r.pool.Query(ctx, query, planID)
	if err != nil {
		r.log.Error("Failed to list tasks by plan ID",
			slog.String("error", err.Error()),
			slog.String("plan_id", planID.String()),
		)
		return nil, fmt.Errorf("repository.ListByPlanID(task): %w", err)
	}
	defer rows.Close()

	tasks := make([]*domain.Task, 0)

	for rows.Next() {
		var m model.TaskModel
		err := rows.Scan(
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
			r.log.Error("Failed to scan task row",
				slog.String("error", err.Error()),
				slog.String("plan_id", planID.String()),
			)
			return nil, err
		}

		tasks = append(tasks, converter.ToTaskEntity(&m))
	}

	if err := rows.Err(); err != nil {
		r.log.Error("Rows iteration error in ListByPlanID",
			slog.String("error", err.Error()),
			slog.String("plan_id", planID.String()),
		)
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) CompleteTestingTask(ctx context.Context, planID uuid.UUID) error {
	query := `
	UPDATE tasks
	SET 
		status = 'done',
		updated_at = NOW()
	WHERE 
		plan_id = $1
	AND 
		title = 'Пройти тестирование'
	`

	result, err := r.pool.Exec(ctx, query, planID)
	if err != nil {
		r.log.Error("Failed to complete testing task",
			slog.String("error", err.Error()),
			slog.String("plan_id", planID.String()),
		)
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

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

var ErrPlanNotFound = errors.New("plan not found")

type planRepository struct {
	pool *pgxpool.Pool
}

func NewPlanRepository(pool *pgxpool.Pool) *planRepository {
	return &planRepository{
		pool: pool,
	}
}

func (r *planRepository) Create(ctx context.Context, plan *domain.Plan) (uuid.UUID, error) {
	var id uuid.UUID

	m := converter.ToPlanModel(plan)

	query := `
		INSERT INTO plans (
			employee_id,
			created_by,
			title,
			description,
			creation_type,
			progress,
			status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		m.EmployeeID,
		m.CreatedBy,
		m.Title,
		m.Description,
		m.CreationType,
		m.Progress,
		m.Status,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("repository.Create(plan): %w", err)
	}

	return id, nil
}

func (r *planRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	var m model.PlanModel
	query := `
		SELECT id, 
			employee_id,
			created_by,
			title,
			description,
			creation_type,
			progress,
			status,
			created_at,
			updated_at

		FROM plans
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.EmployeeID,
		&m.CreatedBy,
		&m.Title,
		&m.Description,
		&m.CreationType,
		&m.Progress,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("repository.GetByID(plan): %w", err)
	}
	return converter.ToPlanEntity(&m), nil
}

func (r *planRepository) RecalculateProgress(ctx context.Context, planID uuid.UUID) error {
	query := `
	UPDATE plans
	SET
	progress = x.progress,
	status =
	CASE
		WHEN x.progress=100
		THEN 'completed'
		ELSE status
	END,
	updated_at=NOW()

	FROM (
		SELECT
		$1::uuid id,
		COALESCE(
		ROUND(
			COUNT(*) FILTER(
				WHERE status='done'
			)
			*100.0/
			NULLIF(COUNT(*),0)
		),0
		)::int progress

		FROM tasks
		WHERE plan_id=$1
	)x

	WHERE plans.id=x.id
	`

	_, err := r.pool.Exec(ctx, query, planID)

	return err
}

func (r *planRepository) ListByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]*domain.Plan, error) {

	query := `
        SELECT
            id,
            employee_id,
            created_by,
            title,
            description,
            creation_type,
            progress,
            status,
            created_at,
            updated_at
        FROM plans
        WHERE employee_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.pool.Query(ctx, query, employeeID)
	if err != nil {
		return nil, fmt.Errorf("repository.ListByEmployeeID(plan): %w", err)
	}
	defer rows.Close()

	plans := make([]*domain.Plan, 0)

	for rows.Next() {

		var m model.PlanModel

		err := rows.Scan(
			&m.ID,
			&m.EmployeeID,
			&m.CreatedBy,
			&m.Title,
			&m.Description,
			&m.CreationType,
			&m.Progress,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		plans = append(plans, converter.ToPlanEntity(&m))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func (r *planRepository) GetEmployeePlan(ctx context.Context, employeeID uuid.UUID, planID uuid.UUID) (*domain.Plan, error) {

	var m model.PlanModel

	query := `
        SELECT
            id,
            employee_id,
            created_by,
            title,
            description,
            creation_type,
            progress,
            status,
            created_at,
            updated_at
        FROM plans
        WHERE id=$1
          AND employee_id=$2
    `

	err := r.pool.QueryRow(ctx, query, planID, employeeID).Scan(
		&m.ID,
		&m.EmployeeID,
		&m.CreatedBy,
		&m.Title,
		&m.Description,
		&m.CreationType,
		&m.Progress,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}

		return nil, err
	}

	return converter.ToPlanEntity(&m), nil
}

func (r *planRepository) ListByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.Plan, error) {
	query := `
		SELECT id, employee_id, created_by, title, description, creation_type, 
		       progress, status, created_at, updated_at
		FROM plans
		WHERE created_by = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, managerID)
	if err != nil {
		return nil, fmt.Errorf("query plans by manager: %w", err)
	}
	defer rows.Close()

	var planModels []*model.PlanModel

	for rows.Next() {
		var m model.PlanModel

		err := rows.Scan(
			&m.ID,
			&m.EmployeeID,
			&m.CreatedBy,
			&m.Title,
			&m.Description,
			&m.CreationType,
			&m.Progress,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}

		planModels = append(planModels, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	plans := make([]*domain.Plan, 0, len(planModels))
	for _, m := range planModels {
		plans = append(plans, converter.ToPlanEntity(m))
	}

	return plans, nil
}

func (r *planRepository) GetByIDWithTasks(ctx context.Context, id uuid.UUID) (*domain.PlanWithTasks, error) {
	plan, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, plan_id, title, description, position, status, created_at, updated_at
		FROM tasks
		WHERE plan_id = $1
		ORDER BY position ASC, created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []*model.TaskModel

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
			return nil, fmt.Errorf("scan task: %w", err)
		}

		taskModels = append(taskModels, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	tasks := make([]*domain.Task, 0, len(taskModels))
	for _, m := range taskModels {
		tasks = append(tasks, converter.ToTaskEntity(m))
	}

	return &domain.PlanWithTasks{
		Plan:  plan,
		Tasks: tasks,
	}, nil
}

func (r *planRepository) Update(ctx context.Context, planID uuid.UUID, title string, description *string) error {
	query := `UPDATE plans SET title = $1, updated_at = NOW()`
	args := []any{title}
	argIndex := 2

	if description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *description)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, planID)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update plan: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrPlanNotFound
	}

	return nil
}

func (r *planRepository) Delete(ctx context.Context, planID uuid.UUID) error {
	query := `
		DELETE FROM plans
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, planID)
	if err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrPlanNotFound
	}

	return nil
}

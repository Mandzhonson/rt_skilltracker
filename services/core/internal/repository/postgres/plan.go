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

func (r *planRepository) CreateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	var id uuid.UUID

	m := converter.ToTaskModel(task)

	query := `
		INSERT INTO tasks (
			plan_id,
			title,
			description,
			position,
			status
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		m.PlanID,
		m.Title,
		m.Description,
		m.Status,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repository.CreateTask(task): %w", err)
	}

	return id, nil
}

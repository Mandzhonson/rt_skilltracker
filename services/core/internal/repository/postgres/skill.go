package postgres

import (
	"context"

	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type skillRepository struct {
	pool *pgxpool.Pool
}

func NewSkillRepository(pool *pgxpool.Pool) *skillRepository {
	return &skillRepository{
		pool: pool,
	}
}

func (r *skillRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Skill, error) {

	query := `
	SELECT
		s.id,
		s.name,
		s.category,
		s.description,
		s.created_at
	FROM user_skills us
	JOIN skills s
		ON s.id = us.skill_id
	WHERE us.user_id = $1
	ORDER BY s.name;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Skill

	for rows.Next() {
		var skill model.SkillModel

		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Category,
			&skill.Description,
			&skill.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, converter.ToSkillEntity(&skill))
	}

	return result, rows.Err()
}

func (r *skillRepository) GetByName(ctx context.Context, name string) (*domain.Skill, error) {

	query := `
	SELECT
		id,
		name,
		category,
		description,
		created_at
	FROM skills
	WHERE name=$1;
	`

	var skill model.SkillModel

	err := r.pool.QueryRow(ctx, query, name).Scan(
		&skill.ID,
		&skill.Name,
		&skill.Category,
		&skill.Description,
		&skill.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return converter.ToSkillEntity(&skill), nil
}

func (r *skillRepository) Create(ctx context.Context, skill *domain.Skill) (uuid.UUID, error) {
	query := `
	INSERT INTO skills(
		name,
		category,
		description
	)
	VALUES($1,$2,$3)
	RETURNING id;
	`

	var id uuid.UUID

	err := r.pool.QueryRow(
		ctx,
		query,
		skill.Name,
		skill.Category,
		skill.Description,
	).Scan(&id)

	return id, err
}

func (r *skillRepository) AttachToUser(ctx context.Context, userID uuid.UUID, skillID uuid.UUID, planID uuid.UUID) error {

	query := `
	INSERT INTO user_skills(
		user_id,
		skill_id,
		plan_id
	)
	VALUES($1,$2,$3)
	ON CONFLICT(user_id, skill_id)
	DO NOTHING;
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		userID,
		skillID,
		planID,
	)

	return err
}

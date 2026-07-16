package postgres

import (
	"context"
	"log/slog"

	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type skillRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewSkillRepository(pool *pgxpool.Pool, log *slog.Logger) *skillRepository {
	return &skillRepository{
		pool: pool,
		log:  log,
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
		r.log.Error("Failed to list skills by user ID",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
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
			r.log.Error("Failed to scan skill row",
				slog.String("error", err.Error()),
				slog.String("user_id", userID.String()),
			)
			return nil, err
		}

		result = append(result, converter.ToSkillEntity(&skill))
	}

	if err = rows.Err(); err != nil {
		r.log.Error("Rows iteration error in ListByUserID",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, err
	}

	return result, nil
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
		r.log.Error("Failed to get skill by name",
			slog.String("error", err.Error()),
			slog.String("name", name),
		)
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

	if err != nil {
		r.log.Error("Failed to create skill",
			slog.String("error", err.Error()),
			slog.String("name", skill.Name),
			slog.String("category", string(skill.Category)),
		)
		return uuid.Nil, err
	}

	return id, nil
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

	if err != nil {
		r.log.Error("Failed to attach skill to user",
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
			slog.String("skill_id", skillID.String()),
			slog.String("plan_id", planID.String()),
		)
	}

	return err
}

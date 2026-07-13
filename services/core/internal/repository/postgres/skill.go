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

func (r *skillRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserSkill, error) {
	query := `
		SELECT
			id,
			user_id,
			plan_id,
			name,
			confirmed_at
		FROM user_skills
		WHERE user_id=$1
		ORDER BY confirmed_at DESC
		`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.UserSkill
	for rows.Next() {
		var skill model.UserSkillModel
		err := rows.Scan(
			&skill.ID,
			&skill.UserID,
			&skill.PlanID,
			&skill.Name,
			&skill.ConfirmedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, converter.ToUserSkillEntity(&skill))
	}
	return result, nil
}

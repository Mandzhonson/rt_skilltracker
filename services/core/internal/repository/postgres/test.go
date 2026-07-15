package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTestNotFound = errors.New("test not found")

type testRepository struct {
	pool *pgxpool.Pool
}

func NewTestRepository(pool *pgxpool.Pool) *testRepository {
	return &testRepository{
		pool: pool,
	}
}

func (r *testRepository) CreateWithQuestions(ctx context.Context, test *domain.Test, questions []*domain.Question) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var testID uuid.UUID

	err = tx.QueryRow(ctx,
		`
		INSERT INTO tests(
			plan_id
		)
		VALUES($1)
		RETURNING id
		`,
		test.PlanID,
	).Scan(&testID)

	if err != nil {
		return uuid.Nil, err
	}

	for i, q := range questions {

		var questionID uuid.UUID

		err = tx.QueryRow(ctx,
			`
			INSERT INTO questions(
				plan_id,
				question_text,
				option_a,
				option_b,
				option_c,
				option_d,
				correct_option,
				ai_generated
			)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8)
			RETURNING id
			`,
			q.PlanID,
			q.QuestionText,
			q.OptionA,
			q.OptionB,
			q.OptionC,
			q.OptionD,
			q.CorrectOption,
			q.AIGenerated,
		).Scan(&questionID)

		if err != nil {
			return uuid.Nil, err
		}

		_, err = tx.Exec(ctx,
			`
			INSERT INTO test_questions(
				test_id,
				question_id,
				position
			)
			VALUES($1,$2,$3)
			`,
			testID,
			questionID,
			i+1,
		)

		if err != nil {
			return uuid.Nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return testID, nil
}

func (r *testRepository) GetByPlanID(ctx context.Context, planID uuid.UUID) (*domain.Test, error) {
	query := `
	SELECT
		id,
		plan_id,
		created_at
	FROM tests
	WHERE plan_id=$1
	`

	var entity model.TestModel

	err := r.pool.QueryRow(ctx, query, planID).Scan(
		&entity.ID,
		&entity.PlanID,
		&entity.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTestNotFound
		}

		return nil, err
	}

	return converter.ToTestEntity(&entity), nil
}

func (r *testRepository) GetQuestions(ctx context.Context, testID uuid.UUID) ([]*domain.Question, error) {
	query := `
	SELECT
		q.id,
		q.plan_id,
		q.question_text,
		q.option_a,
		q.option_b,
		q.option_c,
		q.option_d,
		q.correct_option,
		q.ai_generated,
		q.created_at
	FROM questions q
	JOIN test_questions tq
	ON tq.question_id=q.id
	WHERE tq.test_id=$1
	ORDER BY tq.position
	`

	rows, err := r.pool.Query(ctx, query, testID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []*domain.Question

	for rows.Next() {

		var entity model.QuestionModel

		err := rows.Scan(
			&entity.ID,
			&entity.PlanID,
			&entity.QuestionText,
			&entity.OptionA,
			&entity.OptionB,
			&entity.OptionC,
			&entity.OptionD,
			&entity.CorrectOption,
			&entity.AIGenerated,
			&entity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		questions = append(questions, converter.ToQuestionEntity(&entity))
	}
	return questions, nil
}

func (r *testRepository) CreateAttempt(ctx context.Context, attempt *domain.TestAttempt) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.pool.QueryRow(ctx,
		`
		INSERT INTO test_attempts(
			test_id,
			user_id,
			score,
			total,
			passed
		)
		VALUES($1,$2,$3,$4,$5)
		RETURNING id
		`,
		attempt.TestID,
		attempt.UserID,
		attempt.Score,
		attempt.Total,
		attempt.Passed,
	).Scan(&id)

	return id, err
}

func (r *testRepository) CreateAnswers(ctx context.Context, answers []*domain.TestAnswer) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, a := range answers {
		_, err = tx.Exec(ctx,
			`
			INSERT INTO test_answers(
				attempt_id,
				question_id,
				selected_option,
				is_correct
			)
			VALUES($1,$2,$3,$4)
			`,
			a.AttemptID,
			a.QuestionID,
			a.SelectedOption,
			a.IsCorrect,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *testRepository) GetAttempt(ctx context.Context, testID uuid.UUID, userID uuid.UUID) (*domain.TestAttempt, error) {
	var attempt domain.TestAttempt

	err := r.pool.QueryRow(ctx,
		`
		SELECT
			id,
			test_id,
			user_id,
			score,
			total,
			passed,
			ai_feedback,
			started_at,
			finished_at
		FROM test_attempts
		WHERE test_id=$1
		AND user_id=$2
		ORDER BY started_at DESC
		LIMIT 1
		`,
		testID,
		userID,
	).Scan(
		&attempt.ID,
		&attempt.TestID,
		&attempt.UserID,
		&attempt.Score,
		&attempt.Total,
		&attempt.Passed,
		&attempt.AIFeedback,
		&attempt.StartedAt,
		&attempt.FinishedAt,
	)

	if err != nil {
		return nil, err
	}

	return &attempt, nil
}

func (r *testRepository) FinishAttempt(ctx context.Context, attemptID uuid.UUID, score int, passed bool) error {

	_, err := r.pool.Exec(ctx,
		`
		UPDATE test_attempts
		SET
			score=$1,
			passed=$2,
			finished_at=NOW()
		WHERE id=$3
		`,
		score,
		passed,
		attemptID,
	)

	return err
}

package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/converter"
	"core_service/internal/repository/model"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTestNotFound = errors.New("test not found")

type testRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewTestRepository(pool *pgxpool.Pool, log *slog.Logger) *testRepository {
	return &testRepository{
		pool: pool,
		log:  log,
	}
}

func (r *testRepository) CreateWithQuestions(ctx context.Context, test *domain.Test, questions []*domain.Question) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error("Failed to begin transaction for CreateWithQuestions",
			slog.String("error", err.Error()),
			slog.String("plan_id", test.PlanID.String()),
		)
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
		r.log.Error("Failed to create test",
			slog.String("error", err.Error()),
			slog.String("plan_id", test.PlanID.String()),
		)
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
			r.log.Error("Failed to create question",
				slog.String("error", err.Error()),
				slog.String("plan_id", q.PlanID.String()),
				slog.String("question_text", q.QuestionText),
			)
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
			r.log.Error("Failed to associate question with test",
				slog.String("error", err.Error()),
				slog.String("test_id", testID.String()),
				slog.String("question_id", questionID.String()),
			)
			return uuid.Nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error("Failed to commit transaction for CreateWithQuestions",
			slog.String("error", err.Error()),
			slog.String("plan_id", test.PlanID.String()),
		)
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

		r.log.Error("Failed to get test by plan ID",
			slog.String("error", err.Error()),
			slog.String("plan_id", planID.String()),
		)
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
		r.log.Error("Failed to get questions for test",
			slog.String("error", err.Error()),
			slog.String("test_id", testID.String()),
		)
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
			r.log.Error("Failed to scan question row",
				slog.String("error", err.Error()),
				slog.String("test_id", testID.String()),
			)
			return nil, err
		}
		questions = append(questions, converter.ToQuestionEntity(&entity))
	}

	if err = rows.Err(); err != nil {
		r.log.Error("Rows iteration error in GetQuestions",
			slog.String("error", err.Error()),
			slog.String("test_id", testID.String()),
		)
		return nil, err
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

	if err != nil {
		r.log.Error("Failed to create test attempt",
			slog.String("error", err.Error()),
			slog.String("test_id", attempt.TestID.String()),
			slog.String("user_id", attempt.UserID.String()),
		)
		return uuid.Nil, err
	}

	return id, nil
}

func (r *testRepository) CreateAnswers(ctx context.Context, answers []*domain.TestAnswer) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error("Failed to begin transaction for CreateAnswers",
			slog.String("error", err.Error()),
		)
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
			r.log.Error("Failed to create test answer",
				slog.String("error", err.Error()),
				slog.String("attempt_id", a.AttemptID.String()),
				slog.String("question_id", a.QuestionID.String()),
			)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error("Failed to commit transaction for CreateAnswers",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		r.log.Error("Failed to get test attempt",
			slog.String("error", err.Error()),
			slog.String("test_id", testID.String()),
			slog.String("user_id", userID.String()),
		)
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

	if err != nil {
		r.log.Error("Failed to finish test attempt",
			slog.String("error", err.Error()),
			slog.String("attempt_id", attemptID.String()),
			slog.Int("score", score),
			slog.Bool("passed", passed),
		)
	}

	return err
}


CREATE TABLE tests (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id    UUID        NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    skill_id   UUID        NOT NULL REFERENCES skills(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE test_questions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    test_id     UUID     NOT NULL REFERENCES tests(id)     ON DELETE CASCADE,
    question_id UUID     NOT NULL REFERENCES questions(id),
    position    SMALLINT NOT NULL DEFAULT 0,
    UNIQUE (test_id, question_id)
);

CREATE TABLE test_attempts (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    test_id      UUID        NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
    employee_id  UUID        NOT NULL REFERENCES employees(id),
    score        SMALLINT    NOT NULL DEFAULT 0,  -- кол-во правильных ответов
    total        SMALLINT    NOT NULL DEFAULT 10,
    passed       BOOLEAN     NOT NULL DEFAULT FALSE, -- score/total >= 0.7
    ai_feedback  TEXT,                              -- обратная связь от AI после теста
    started_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at  TIMESTAMPTZ
);

CREATE TABLE test_answers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id  UUID     NOT NULL REFERENCES test_attempts(id) ON DELETE CASCADE,
    question_id UUID     NOT NULL REFERENCES questions(id),
    selected_option CHAR(1) NOT NULL CHECK (selected_option IN ('A', 'B', 'C', 'D')),
    is_correct  BOOLEAN  NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_tests_plan_id            ON tests(plan_id);
CREATE INDEX idx_test_questions_test_id   ON test_questions(test_id);
CREATE INDEX idx_test_attempts_test_id    ON test_attempts(test_id);
CREATE INDEX idx_test_attempts_employee   ON test_attempts(employee_id);
CREATE INDEX idx_test_answers_attempt_id  ON test_answers(attempt_id);
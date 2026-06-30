CREATE TABLE questions (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    skill_id       UUID         NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    question_text  TEXT         NOT NULL,
    option_a       VARCHAR(500) NOT NULL,
    option_b       VARCHAR(500) NOT NULL,
    option_c       VARCHAR(500) NOT NULL,
    option_d       VARCHAR(500) NOT NULL,
    correct_option CHAR(1)      NOT NULL CHECK (correct_option IN ('A', 'B', 'C', 'D')),
    ai_generated   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_skill ON questions(skill_id);
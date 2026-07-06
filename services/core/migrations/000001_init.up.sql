CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    avatar_key TEXT,
    
    role VARCHAR(50) NOT NULL DEFAULT 'employee'
        CHECK (role IN ('employee', 'manager', 'admin')),

    manager_id UUID REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_manager_id ON users(manager_id);



CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    jti VARCHAR(255) UNIQUE NOT NULL,

    token_hash TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,

    revoked BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);



CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name VARCHAR(255) NOT NULL UNIQUE,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



CREATE TABLE user_skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,

    confirmed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, skill_id)
);

CREATE INDEX idx_user_skills_user ON user_skills(user_id);
CREATE INDEX idx_user_skills_skill ON user_skills(skill_id);



CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    created_by UUID NOT NULL REFERENCES users(id),

    skill_id UUID NOT NULL REFERENCES skills(id),

    title VARCHAR(255) NOT NULL,

    description TEXT,

    status VARCHAR(50) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'completed', 'archived')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_plans_user_id ON plans(user_id);
CREATE INDEX idx_plans_skill_id ON plans(skill_id);
CREATE INDEX idx_plans_status ON plans(status);



CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    position SMALLINT NOT NULL DEFAULT 0,

    status VARCHAR(50) NOT NULL DEFAULT 'todo'
        CHECK (status IN ('todo', 'in_progress', 'done')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_plan_id ON tasks(plan_id);
CREATE INDEX idx_tasks_status ON tasks(status);



CREATE TABLE task_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,

    author_id UUID NOT NULL REFERENCES users(id),

    message TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_task_comments_task ON task_comments(task_id);



CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,

    question_text TEXT NOT NULL,

    option_a VARCHAR(500) NOT NULL,

    option_b VARCHAR(500) NOT NULL,

    option_c VARCHAR(500) NOT NULL,

    option_d VARCHAR(500) NOT NULL,

    correct_option CHAR(1) NOT NULL
        CHECK (correct_option IN ('A', 'B', 'C', 'D')),

    ai_generated BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_skill ON questions(skill_id);



CREATE TABLE tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,

    skill_id UUID NOT NULL REFERENCES skills(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tests_plan_id ON tests(plan_id);



CREATE TABLE test_questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,

    question_id UUID NOT NULL REFERENCES questions(id),

    position SMALLINT NOT NULL DEFAULT 0,

    UNIQUE(test_id, question_id)
);

CREATE INDEX idx_test_questions_test_id ON test_questions(test_id);



CREATE TABLE test_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,

    user_id UUID NOT NULL REFERENCES users(id),

    score SMALLINT NOT NULL DEFAULT 0,

    total SMALLINT NOT NULL DEFAULT 10,

    passed BOOLEAN NOT NULL DEFAULT FALSE,

    ai_feedback TEXT,

    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    finished_at TIMESTAMPTZ
);

CREATE INDEX idx_test_attempts_test_id ON test_attempts(test_id);
CREATE INDEX idx_test_attempts_user_id ON test_attempts(user_id);



CREATE TABLE test_answers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    attempt_id UUID NOT NULL REFERENCES test_attempts(id) ON DELETE CASCADE,

    question_id UUID NOT NULL REFERENCES questions(id),

    selected_option CHAR(1) NOT NULL
        CHECK (selected_option IN ('A', 'B', 'C', 'D')),

    is_correct BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_test_answers_attempt_id ON test_answers(attempt_id);
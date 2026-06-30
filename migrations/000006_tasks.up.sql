CREATE TABLE tasks (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id      UUID         NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    position     SMALLINT     NOT NULL DEFAULT 0, -- порядок шага в плане
    status       VARCHAR(50)  NOT NULL DEFAULT 'todo'
                              CHECK (status IN ('todo', 'in_progress', 'done')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE task_comments (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id    UUID        NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    author_id  UUID        NOT NULL REFERENCES employees(id), -- работодатель
    message    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_plan_id      ON tasks(plan_id);
CREATE INDEX idx_tasks_status       ON tasks(status);
CREATE INDEX idx_task_comments_task ON task_comments(task_id);
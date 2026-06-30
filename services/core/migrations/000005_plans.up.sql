CREATE TABLE plans (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id  UUID         NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    created_by   UUID         NOT NULL REFERENCES employees(id), -- работодатель
    skill_id     UUID         NOT NULL REFERENCES skills(id),
    title        VARCHAR(255) NOT NULL,        -- сгенерировано AI
    description  TEXT,                         -- сгенерировано AI
    status       VARCHAR(50)  NOT NULL DEFAULT 'active'
                              CHECK (status IN ('active', 'completed', 'archived')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_plans_employee_id ON plans(employee_id);
CREATE INDEX idx_plans_skill_id    ON plans(skill_id);
CREATE INDEX idx_plans_status      ON plans(status);
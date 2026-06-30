CREATE TABLE skills (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE employee_skills (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID        NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    skill_id    UUID        NOT NULL REFERENCES skills(id)    ON DELETE CASCADE,
    confirmed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (employee_id, skill_id)
);

CREATE INDEX idx_employee_skills_employee ON employee_skills(employee_id);
CREATE INDEX idx_employee_skills_skill    ON employee_skills(skill_id);
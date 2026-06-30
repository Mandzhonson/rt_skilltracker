CREATE TABLE employees (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auth_id       UUID         NOT NULL UNIQUE REFERENCES credentials(id) ON DELETE CASCADE,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    manager_id    UUID         REFERENCES employees(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_employees_auth_id    ON employees(auth_id);
CREATE INDEX idx_employees_manager_id ON employees(manager_id);
CREATE TABLE credentials (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'employee'
                               CHECK (role IN ('employee', 'manager', 'admin')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    credential_id UUID        NOT NULL REFERENCES credentials(id) ON DELETE CASCADE,
    token_hash    TEXT        NOT NULL UNIQUE,
    expires_at    TIMESTAMPTZ NOT NULL,
    revoked       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE blacklisted_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jti        VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_credential_id ON refresh_tokens(credential_id);
CREATE INDEX idx_refresh_tokens_expires_at    ON refresh_tokens(expires_at);
CREATE INDEX idx_blacklisted_tokens_jti       ON blacklisted_tokens(jti);
CREATE TABLE
    refresh_tokens (
        jti UUID PRIMARY KEY,
        user_id BIGINT NOT NULL,
        expires_at TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        revoked_at TIMESTAMPTZ
    );

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
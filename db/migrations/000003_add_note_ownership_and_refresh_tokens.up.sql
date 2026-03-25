ALTER TABLE app_notes
ADD COLUMN IF NOT EXISTS user_id BIGINT REFERENCES app_users (id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_app_notes_user_id_created_at ON app_notes (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS app_refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES app_users (id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_app_refresh_tokens_user_id ON app_refresh_tokens (user_id, created_at DESC);

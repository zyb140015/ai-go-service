DROP INDEX IF EXISTS idx_app_refresh_tokens_user_id;
DROP TABLE IF EXISTS app_refresh_tokens;

DROP INDEX IF EXISTS idx_app_notes_user_id_created_at;
ALTER TABLE app_notes DROP COLUMN IF EXISTS user_id;

-- name: CreateRefreshToken :one
INSERT INTO app_refresh_tokens (
    user_id,
    token_hash,
    expires_at
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at;

-- name: GetActiveRefreshToken :one
SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
FROM app_refresh_tokens
WHERE user_id = $1
  AND token_hash = $2
  AND revoked_at IS NULL
  AND expires_at > $3;

-- name: RevokeRefreshToken :execrows
UPDATE app_refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = $1
  AND revoked_at IS NULL;

-- name: RevokeRefreshTokensByUser :execrows
UPDATE app_refresh_tokens
SET revoked_at = NOW()
WHERE user_id = $1
  AND revoked_at IS NULL;

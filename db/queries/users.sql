-- name: CreateUser :one
INSERT INTO app_users (
    email,
    display_name,
    password_hash
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, email, display_name, password_hash, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, display_name, password_hash, created_at, updated_at
FROM app_users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, display_name, password_hash, created_at, updated_at
FROM app_users
WHERE id = $1;

-- name: UpdateUserPassword :one
UPDATE app_users
SET password_hash = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, display_name, password_hash, created_at, updated_at;

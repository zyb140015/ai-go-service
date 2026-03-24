-- name: CreateNote :one
INSERT INTO app_notes (
    title,
    body
) VALUES (
    $1,
    $2
)
RETURNING id, title, body, created_at, updated_at;

-- name: ListNotes :many
SELECT id, title, body, created_at, updated_at
FROM app_notes
ORDER BY created_at DESC, id DESC;

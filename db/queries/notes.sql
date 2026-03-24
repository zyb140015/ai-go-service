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

-- name: GetNoteByID :one
SELECT id, title, body, created_at, updated_at
FROM app_notes
WHERE id = $1;

-- name: UpdateNote :one
UPDATE app_notes
SET title = $2,
    body = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, body, created_at, updated_at;

-- name: DeleteNote :execrows
DELETE FROM app_notes
WHERE id = $1;

-- name: CreateNote :one
INSERT INTO app_notes (
    user_id,
    title,
    body
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, user_id, title, body, created_at, updated_at;

-- name: ListNotes :many
SELECT id, user_id, title, body, created_at, updated_at
FROM app_notes
WHERE user_id = sqlc.arg(user_id)
  AND CASE
    WHEN sqlc.narg(query_text)::text IS NULL OR sqlc.narg(query_text)::text = '' THEN TRUE
    ELSE title ILIKE '%' || sqlc.narg(query_text)::text || '%'
END
ORDER BY
    CASE WHEN sqlc.arg(sort_field)::text = 'title' AND sqlc.arg(sort_direction)::text = 'asc' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_field)::text = 'title' AND sqlc.arg(sort_direction)::text = 'desc' THEN title END DESC,
    CASE WHEN sqlc.arg(sort_field)::text = 'created_at' AND sqlc.arg(sort_direction)::text = 'asc' THEN created_at END ASC,
    CASE WHEN sqlc.arg(sort_field)::text = 'created_at' AND sqlc.arg(sort_direction)::text = 'desc' THEN created_at END DESC,
    id DESC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);

-- name: CountNotes :one
SELECT COUNT(*)
FROM app_notes
WHERE user_id = sqlc.arg(user_id)
  AND CASE
    WHEN sqlc.narg(query_text)::text IS NULL OR sqlc.narg(query_text)::text = '' THEN TRUE
    ELSE title ILIKE '%' || sqlc.narg(query_text)::text || '%'
END;

-- name: GetNoteByID :one
SELECT id, user_id, title, body, created_at, updated_at
FROM app_notes
WHERE id = $1
  AND user_id = $2;

-- name: UpdateNote :one
UPDATE app_notes
SET title = $2,
    body = $3,
    updated_at = NOW()
WHERE id = $1
  AND user_id = $4
RETURNING id, user_id, title, body, created_at, updated_at;

-- name: DeleteNote :execrows
DELETE FROM app_notes
WHERE id = $1
  AND user_id = $2;

-- name: GetChats :many
SELECT * FROM chats;

-- name: CreateChat :one
INSERT INTO chats (id, user_id, title, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

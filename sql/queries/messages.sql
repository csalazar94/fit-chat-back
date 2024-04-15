-- name: CreateMessage :one
INSERT INTO messages (id, chat_id, author_role_id, content, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetMessagesByChatId :many
SELECT * FROM messages WHERE chat_id = $1;

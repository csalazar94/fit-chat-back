-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (id, full_name, email, password, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: CreateUser :one

INSERT INTO users (
    name,
    email,
    password
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
    name    = $2,
    email    = $3,
    email_verified = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;
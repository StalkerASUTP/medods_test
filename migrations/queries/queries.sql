-- name: CreateUser :one
INSERT INTO users (
    id,
    refresh_token_hash,
    refresh_token_expires_at,
    user_agent,
    ip_address,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateToken :one
UPDATE users
SET refresh_token_hash = $2,
    refresh_token_expires_at = $3
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = FALSE
WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;




-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, is_active, created_at;

-- name: EmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1) AS exists;

-- name: UsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) AS exists;

-- name: GetUserByEmail :one 
SELECT id, username, email, password_hash, is_active
FROM users
WHERE email = $1;

-- name: ActivateUser :exec
UPDATE users
SET is_active = TRUE, updated_at = NOW()
WHERE id = $1;
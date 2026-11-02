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

-- name: GetPasswdHash :one
SELECT password_hash
FROM users
WHERE email = $1;

-- профиль

-- name: GetProfile :one
SELECT username, email, COALESCE(title, '')::text AS title, COALESCE(bio, '')::text AS bio
FROM users
WHERE id = $1;

-- name: UpdateProfile :one
UPDATE users
SET username = sqlc.arg(username),
    email = sqlc.arg(email),
    title = NULLIF(sqlc.arg(title)::text, ''),
    bio = NULLIF(sqlc.arg(bio)::text, ''),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING username, email, COALESCE(title, '')::text AS title, COALESCE(bio, '')::text AS bio;

-- привязки соцсетей

-- name: GetUserByIdentity :one
SELECT u.id, u.username, u.is_active
FROM users u
JOIN user_identities i ON i.user_id = u.id
WHERE i.provider = $1 AND i.provider_user_id = $2;

-- name: ListIdentities :many
SELECT provider FROM user_identities WHERE user_id = $1;

-- Привязать или перепривязать (замена без удаления)
-- name: UpsertIdentity :exec
INSERT INTO user_identities (user_id, provider, provider_user_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, provider) DO UPDATE
SET provider_user_id = EXCLUDED.provider_user_id, created_at = NOW();

-- name: DeleteIdentity :execrows
DELETE FROM user_identities WHERE user_id = $1 AND provider = $2;
-- name: UpsertActivationCode :exec
INSERT INTO activation_codes (user_id, code_hash, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (user_id)
DO UPDATE SET code_hash = EXCLUDED.code_hash, expires_at = EXCLUDED.expires_at;

-- name: GetActivationCodeByUserID :one
SELECT user_id, code_hash, expires_at
FROM activation_codes
WHERE user_id = $1;

-- name: DeleteActivationCode :exec
DELETE FROM activation_codes 
WHERE user_id = $1;

-- name: GetPendingEmail :one
SELECT email FROM users
WHERE id = $1 AND is_active = FALSE;

-- name: ActivateUserAndDeleteCode :exec
WITH deleted AS (
    DELETE FROM activation_codes WHERE user_id = $1
)
UPDATE users
SET is_active = TRUE, updated_at = NOW()
WHERE users.id = $1;
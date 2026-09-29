-- name: GetLinkByShareCode :one
SELECT id, created_by, share_code, redirect_timer, redirect_to, valid_until, allowed_redirects, redirects, date_created, only_unique_redirects, is_active
FROM links
WHERE share_code = $1;

-- name: SaveLink :one
INSERT INTO links (created_by, share_code, redirect_timer, redirect_to, valid_until, allowed_redirects, only_unique_redirects, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, share_code, redirect_timer, redirect_to, valid_until, is_active, date_created;

-- name: GetOnlyLinkByShareCode :one
SELECT redirect_to
FROM links
WHERE share_code = $1;

-- name: IncrementRedirectCount :one
UPDATE links 
SET redirects = redirects + 1
WHERE id = $1 AND (allowed_redirects IS NULL OR redirects < allowed_redirects)
RETURNING redirects;

-- name: ConfirmRedirect :one
UPDATE links
SET redirects = redirects + 1
WHERE id = $1
    AND is_active = TRUE
    AND (valid_until IS NULL OR valid_until > NOW())
    AND (allowed_redirects IS NULL OR redirect < allowed_redirects)
RETURNING redirect_to, redirects;

-- name: HasVisited :one 
SELECT EXISTS(
    SELECT 1 FROM link_visits WHERE link_id = $1 AND ip_hash = $2
) AS visited;

-- name: RecordVisit :exec
INSERT INTO link_visits (link_id, ip_hash)
VALUES ($1, $2)
ON CONFLICT (link_id, ip_hash) DO NOTHING;

-- name: SetLinkActive :exec 
UPDATE LINKS
SET is_active = $2
WHERE id = $1 AND created_by = $3;

-- name: DeleteLink :exec
DELETE FROM links 
WHERE id = $1 AND created_by = $2;


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
WITH confirmed AS (
    UPDATE links
    SET redirects = redirects + 1
    WHERE id = sqlc.arg(link_id)
        AND is_active = TRUE
        AND (valid_until IS NULL OR valid_until > NOW())
        AND (allowed_redirects IS NULL OR redirects < allowed_redirects)
    RETURNING id, redirect_to, redirects
), daily AS (
    INSERT INTO link_daily_stats (link_id, day, redirects)
    SELECT id, (NOW() AT TIME ZONE sqlc.arg(tz)::text)::date, 1
    FROM confirmed
    ON CONFLICT (link_id, day)
    DO UPDATE SET redirects = link_daily_stats.redirects + 1
)
SELECT redirect_to, redirects FROM confirmed;

-- name: GetRedirectStats :one
SELECT
    COALESCE((SELECT SUM(redirects) FROM links), 0)::bigint AS total,
    COALESCE((SELECT SUM(redirects) FROM links
              WHERE created_by = sqlc.arg(user_id)::int), 0)::bigint AS total_mine,
    COALESCE((SELECT SUM(redirects) FROM link_daily_stats
              WHERE day = (NOW() AT TIME ZONE sqlc.arg(tz)::text)::date), 0)::bigint AS today,
    COALESCE((SELECT SUM(s.redirects) FROM link_daily_stats s
              JOIN links l ON l.id = s.link_id
              WHERE l.created_by = sqlc.arg(user_id)::int
                AND s.day = (NOW() AT TIME ZONE sqlc.arg(tz)::text)::date), 0)::bigint AS today_mine;

-- name: HasVisited :one 
SELECT EXISTS(
    SELECT 1 FROM link_visits WHERE link_id = $1 AND ip_hash = $2
) AS visited;

-- name: RecordVisit :execrows
INSERT INTO link_visits (link_id, ip_hash)
VALUES ($1, $2)
ON CONFLICT (link_id, ip_hash) DO NOTHING;

-- name: SetLinkActive :exec 
UPDATE LINKS
SET is_active = $3
WHERE id = $1 AND created_by = $2;

-- name: DeleteLink :exec
DELETE FROM links 
WHERE id = $1 AND created_by = $2;
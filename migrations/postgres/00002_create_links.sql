-- +goose Up
CREATE TABLE links (
    id PRIMARY KEY,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    share_code VARCHAR(10) UNIQUE NOT NULL,
    redirect_timer INT NOT NULL DEFAULT 5 CHECK (redirect_timer IN (0,3,5,10,15)),
    redirect_to TEXT NOT NULL,
    valid_until TIMESTAMPTZ,
    allowed_redirects INT CHECK (allowed_redirects NOT NULL OR allowed_redirects >= 1),
    redirects INT NOT NULL DEFAULT 0,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    only_unique_redirects BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_links_created_by ON links (created_by);
CREATE INDEX idx_links_valid_until ON links (valid_until);

-- +goose Down
DROP TABLE links;
-- +goose Up
CREATE TABLE link_visits (
    id SERIAL PRIMARY KEY,
    link_id INT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    ip_hash TEXT NOT NULL,
    visited_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_link_visits_unique ON link_visits (link_id, ip_hash);
CREATE INDEX idx_link_visits_link_id ON link_visits (link_id);

-- +goose Down
DROP TABLE link_visits;

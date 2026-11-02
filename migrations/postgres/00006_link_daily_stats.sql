-- +goose Up
CREATE TABLE link_daily_stats (
    link_id INT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    day DATE NOT NULL,
    redirects INT NOT NULL DEFAULT 0,
    PRIMARY KEY (link_id, day)
);
CREATE INDEX idx_link_daily_stats_day ON link_daily_stats (day);

-- +goose Down
DROP TABLE link_daily_stats;
-- +goose Up
CREATE TABLE activation_codes (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() 
);

CREATE INDEX idx_activation_codes_expire_at ON activation_codes (expires_at);

-- +goose Down
DROP TABLE activation_codes;

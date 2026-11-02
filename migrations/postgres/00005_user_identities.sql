-- +goose Up
CREATE TABLE user_identities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(20) NOT NULL,        
    provider_user_id VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider, provider_user_id),   
    UNIQUE (user_id, provider)            
);
-- +goose Down
DROP TABLE user_identities;
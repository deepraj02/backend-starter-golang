-- +goose Up
ALTER TABLE users ADD COLUMN provider VARCHAR(20) DEFAULT 'email' NOT NULL;

-- Create index for provider field
CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider);

-- Update existing users to have 'email' provider
UPDATE users SET provider = 'email' WHERE provider IS NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN provider;
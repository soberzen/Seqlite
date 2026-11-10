-- +goose Up
CREATE TABLE users(
  id BIGSERIAL PRIMARY KEY,
  user_name VARCHAR(32) NOT NULL,
  email VARCHAR(128) NOT NULL,
  password VARCHAR(255) NOT NULL,
  avatar_url VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email ON users (email)
  WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS users

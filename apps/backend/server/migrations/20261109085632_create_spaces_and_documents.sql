-- +goose Up
CREATE TABLE spaces(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  emoji VARCHAR(32),
  description TEXT NOT NULL DEFAULT '',
  status SMALLINT NOT NULL DEFAULT 0,
  owner_id BIGINT NOT NULL,
  expire_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_spaces_status CHECK (status IN (0, 1)),
  CONSTRAINT ck_spaces_expire_at CHECK ((status = 0 AND expire_at IS NULL) OR (status = 1 AND expire_at IS NOT NULL)),
  CONSTRAINT fk_spaces_owner FOREIGN KEY (owner_id)
    REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE INDEX idx_spaces_owner_id ON spaces (owner_id);

CREATE TABLE documents(
  id BIGSERIAL PRIMARY KEY,
  space_id BIGINT NOT NULL,
  title VARCHAR(255) NOT NULL,
  emoji VARCHAR(32),
  status SMALLINT NOT NULL DEFAULT 0,
  owner_id BIGINT NOT NULL,
  expire_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_documents_status CHECK (status IN (0, 1)),
  CONSTRAINT ck_documents_expire_at CHECK ((status = 0 AND expire_at IS NULL) OR (status = 1 AND expire_at IS NOT NULL)),
  CONSTRAINT fk_documents_space FOREIGN KEY (space_id)
    REFERENCES spaces (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_documents_owner FOREIGN KEY (owner_id)
    REFERENCES users (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE INDEX idx_documents_space_id ON documents (space_id);
CREATE INDEX idx_documents_owner_id ON documents (owner_id);

-- +goose Down
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS spaces;

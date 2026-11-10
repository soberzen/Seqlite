-- +goose Up
CREATE TABLE space_grants(
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  space_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  granted_by BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_space_grants_space
    FOREIGN KEY (space_id) REFERENCES spaces (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_space_grants_role
    FOREIGN KEY (role_id) REFERENCES roles (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_space_grants_granted_by
    FOREIGN KEY (granted_by) REFERENCES users (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE UNIQUE INDEX uq_space_grants_user_space
    ON space_grants (user_id, space_id);
CREATE INDEX idx_space_grants_space_id ON space_grants (space_id);
CREATE INDEX idx_space_grants_role_id ON space_grants (role_id);
CREATE INDEX idx_space_grants_granted_by ON space_grants (granted_by);

CREATE TABLE document_grants (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  document_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  granted_by BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_document_grants_user
    FOREIGN KEY (user_id) REFERENCES users (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_document_grants_document
    FOREIGN KEY (document_id) REFERENCES documents (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_document_grants_role
    FOREIGN KEY (role_id) REFERENCES roles (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_document_grants_granted_by
    FOREIGN KEY (granted_by) REFERENCES users (id)
    ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE UNIQUE INDEX uq_document_grants_user_document
    ON document_grants (user_id, document_id);
CREATE INDEX idx_document_grants_document_id ON document_grants (document_id);
CREATE INDEX idx_document_grants_role_id ON document_grants (role_id);
CREATE INDEX idx_document_grants_granted_by ON document_grants (granted_by);
-- +goose Down
DROP TABLE IF EXISTS document_grants;
DROP TABLE IF EXISTS space_grants;

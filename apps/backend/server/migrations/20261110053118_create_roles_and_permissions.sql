-- +goose Up
CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(32) NOT NULL,
  name VARCHAR(64) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_roles_code ON roles (code);

CREATE TABLE document_permissions (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(128) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_document_permissions_code ON document_permissions (code);

CREATE TABLE space_permissions (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(128) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_space_permissions_code ON space_permissions (code);

CREATE TABLE role_document_permissions (
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  PRIMARY KEY (role_id, permission_id),
  CONSTRAINT fk_role_document_permissions_role FOREIGN KEY (role_id)
    REFERENCES roles (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_role_document_permissions_permission FOREIGN KEY (permission_id)
    REFERENCES document_permissions (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE INDEX idx_role_document_permissions_permission_id
  ON role_document_permissions (permission_id);

CREATE TABLE role_space_permissions (
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  PRIMARY KEY (role_id, permission_id),
  CONSTRAINT fk_role_space_permissions_role FOREIGN KEY (role_id)
    REFERENCES roles (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_role_space_permissions_permission FOREIGN KEY (permission_id)
    REFERENCES space_permissions (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE INDEX idx_role_space_permissions_permission_id
  ON role_space_permissions (permission_id);

-- +goose Down
DROP TABLE IF EXISTS role_space_permissions;
DROP TABLE IF EXISTS  role_document_permissions;
DROP TABLE IF EXISTS space_permissions;
DROP TABLE IF EXISTS document_permissions;
DROP TABLE IF EXISTS roles;

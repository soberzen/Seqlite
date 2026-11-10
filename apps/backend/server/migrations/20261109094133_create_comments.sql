-- +goose Up
CREATE TABLE comments(
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  document_id BIGINT NOT NULL,
  parent_id BIGINT,
  content TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,

  CONSTRAINT uq_comments_id_document_id UNIQUE (id, document_id),
  CONSTRAINT ck_comments_not_self_parent
    CHECK (parent_id IS NULL OR parent_id <> id),
  CONSTRAINT chk_comments_content_not_empty
    CHECK (length(trim(content)) > 0),
  CONSTRAINT fk_comments_user FOREIGN KEY (user_id)
    REFERENCES users (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_comments_document FOREIGN KEY (document_id)
    REFERENCES documents (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT fk_comments_parent_same_document
    FOREIGN KEY (parent_id, document_id)
    REFERENCES comments (id, document_id)
    ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE INDEX idx_comments_user_id ON comments (user_id);
CREATE INDEX idx_comments_document_active ON comments (document_id, created_at)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_comments_parent_id ON comments (parent_id)
  WHERE parent_id IS NOT NULL;
CREATE INDEX idx_comments_deleted_at ON comments (deleted_at)
  WHERE deleted_at IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS comments;

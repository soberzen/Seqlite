-- +goose Up
ALTER TABLE documents
  ADD COLUMN yjs_revision BIGINT NOT NULL DEFAULT 0;

CREATE TABLE document_revisions (
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    revision    BIGINT NOT NULL CHECK (revision > 0),
    yjs_update  BYTEA NOT NULL CHECK (octet_length(yjs_update) > 0),
    actor_id    BIGINT NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (document_id, revision)
);

CREATE TABLE document_checkpoints (
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    revision    BIGINT NOT NULL CHECK (revision > 0),
    yjs_state   BYTEA NOT NULL CHECK (octet_length(yjs_state) > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (document_id, revision)
);

-- +goose Down
ALTER TABLE documents DROP COLUMN yjs_revision;
DROP TABLE document_checkpoints;
DROP TABLE document_revisions;

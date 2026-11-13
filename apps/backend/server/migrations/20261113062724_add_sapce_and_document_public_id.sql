-- +goose Up
ALTER TABLE spaces
  ADD COLUMN public_id UUID NOT NULL DEFAULT gen_random_uuid(),
  ADD CONSTRAINT uq_spaces_public_id UNIQUE (public_id);

ALTER TABLE documents
    ADD COLUMN public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    ADD CONSTRAINT uq_documents_public_id UNIQUE (public_id);
-- +goose Down
ALTER TABLE documents DROP COLUMN public_id;
ALTER TABLE spaces DROP COLUMN public_id;

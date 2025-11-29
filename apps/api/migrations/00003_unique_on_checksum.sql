-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD CONSTRAINT unique_checksum UNIQUE (company_id, checksum);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP CONSTRAINT unique_checksum;
-- +goose StatementEnd

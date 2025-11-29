-- +goose Up
-- +goose StatementBegin
ALTER TABLE company_users ADD COLUMN added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE company_users DROP COLUMN added_at;
-- +goose StatementEnd


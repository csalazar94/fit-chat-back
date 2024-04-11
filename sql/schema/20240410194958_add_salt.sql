-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN encoded_hash TEXT NOT NULL DEFAULT 'encoded_hash';
ALTER TABLE users DROP COLUMN password;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN encoded_hash;
ALTER TABLE users ADD COLUMN password TEXT NOT NULL DEFAULT 'password';
-- +goose StatementEnd

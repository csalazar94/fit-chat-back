-- +goose Up
-- +goose StatementBegin
-- system
-- assistant
-- tool
-- user
INSERT INTO author_roles (id, name, created_at, updated_at)
VALUES
    (1, 'Sistema', timezone('utc', now()), timezone('utc', now())),
    (2, 'Asistente', timezone('utc', now()), timezone('utc', now())),
    (3, 'Herramienta', timezone('utc', now()), timezone('utc', now())),
    (4, 'Usuario', timezone('utc', now()), timezone('utc', now()));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM author_roles
WHERE id IN (1, 2, 3, 4);
-- +goose StatementEnd

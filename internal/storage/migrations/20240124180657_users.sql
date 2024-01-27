-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id           INTEGER PRIMARY KEY,
    mailing_time INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS postcards (
    Name VARCHAR(255) NOT NULL,
    path TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS postcards;
-- +goose StatementEnd

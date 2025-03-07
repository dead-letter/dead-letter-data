-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rider_ (
    id_ uuid PRIMARY KEY REFERENCES user_(id_) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rider_;
-- +goose StatementEnd

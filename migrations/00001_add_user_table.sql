-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS user_ (
    id_ uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at_ TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version_ integer NOT NULL DEFAULT 1,
    email_ CITEXT UNIQUE NOT NULL,
	password_hash_ BYTEA NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_;
-- +goose StatementEnd

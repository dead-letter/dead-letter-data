-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS user_ (
    id_ uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at_ TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    email_ CITEXT UNIQUE NOT NULL,
    version_ integer NOT NULL DEFAULT 1
);

INSERT INTO user_ (email_, password_hash_) VALUES
	('test@example.com', crypt('password', gen_salt('bf')))
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_;
-- +goose StatementEnd

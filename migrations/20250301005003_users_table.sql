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

INSERT INTO user_ (email_, password_hash_) 
SELECT 'test@example.com', crypt('super_secure_password', gen_salt('bf'))
WHERE NOT EXISTS (
    SELECT 1 FROM user_ WHERE email_ = 'test@example.com'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_;
-- +goose StatementEnd

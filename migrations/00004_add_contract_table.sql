-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS contract_ (
	id_ uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	created_at_ TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	expiry_ TIMESTAMPTZ,
	rider_id_ uuid REFERENCES rider_(id_) ON DELETE CASCADE,
	vendor_id_ uuid REFERENCES vendor_(id_) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contract_;
-- +goose StatementEnd

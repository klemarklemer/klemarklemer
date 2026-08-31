-- +goose Up
-- +goose StatementBegin
ALTER TABLE claims
ADD COLUMN IF NOT EXISTS fraud_signal VARCHAR(512);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE claims
DROP COLUMN IF EXISTS fraud_signal;
-- +goose StatementEnd

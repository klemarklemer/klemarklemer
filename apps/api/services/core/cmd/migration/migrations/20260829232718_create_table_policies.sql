-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS policies (
	"id" SERIAL PRIMARY KEY,
	"policy_number" VARCHAR(64) UNIQUE NOT NULL,
	"policy_holder_name" VARCHAR(255) NOT NULL,
	"vehicle_plate" VARCHAR(32) NOT NULL,
	"vehicle_model" VARCHAR(128) NOT NULL,
	"coverage_type" VARCHAR(64) NOT NULL,
	"max_coverage_amount" DECIMAL(15,2) NOT NULL,
	"deductible_amount" DECIMAL(15,2) NOT NULL,
	"effective_date" TIMESTAMPTZ(6) NOT NULL,
	"expiry_date" TIMESTAMPTZ(6) NOT NULL,
	"status" VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
	"created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS policies;
-- +goose StatementEnd

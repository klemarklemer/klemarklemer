-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS claims_officers (
	"id" SERIAL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"email" VARCHAR(255) UNIQUE NOT NULL,
	"role" VARCHAR(64) NOT NULL DEFAULT 'Claims Officer',
	"current_workload" INT NOT NULL DEFAULT 0,
	"motor_skill_rating" DECIMAL(3,2) NOT NULL DEFAULT 4.0,
	"is_available" BOOLEAN NOT NULL DEFAULT TRUE,
	"created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS claims_officers;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS claims (
	"id" SERIAL PRIMARY KEY,
	"claim_number" VARCHAR(64) UNIQUE NOT NULL,
	"policy_id" INT NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
	"stage" VARCHAR(64) NOT NULL DEFAULT 'INTAKE',
	"document_completeness" VARCHAR(32) NOT NULL DEFAULT 'INCOMPLETE',
	"survey_required" BOOLEAN NOT NULL DEFAULT FALSE,
	"claim_type" VARCHAR(64) NOT NULL DEFAULT 'MOTOR',
	"severity" VARCHAR(32) NOT NULL DEFAULT 'MEDIUM',
	"incident_description" TEXT,
	"estimated_loss" DECIMAL(15,2) NOT NULL DEFAULT 0,
	"approved_amount" DECIMAL(15,2) NOT NULL DEFAULT 0,
	"current_officer_id" INT REFERENCES claims_officers(id) ON DELETE SET NULL,
	"claim_sla_due_at" TIMESTAMPTZ(6),
	"stage_sla_due_at" TIMESTAMPTZ(6),
	"status" VARCHAR(32) NOT NULL DEFAULT 'OPEN',
	"created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS claim_documents (
	"id" SERIAL PRIMARY KEY,
	"claim_id" INT NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
	"document_type" VARCHAR(64) NOT NULL,
	"file_name" VARCHAR(255) NOT NULL,
	"file_url" TEXT NOT NULL,
	"status" VARCHAR(32) NOT NULL DEFAULT 'VERIFIED',
	"extracted_data" TEXT,
	"uploaded_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS claim_events (
	"id" SERIAL PRIMARY KEY,
	"claim_id" INT NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
	"actor_name" VARCHAR(128) NOT NULL,
	"actor_type" VARCHAR(32) NOT NULL,
	"action" VARCHAR(128) NOT NULL,
	"previous_stage" VARCHAR(64),
	"new_stage" VARCHAR(64),
	"payload" TEXT,
	"created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assignments (
	"id" SERIAL PRIMARY KEY,
	"claim_id" INT UNIQUE NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
	"officer_id" INT NOT NULL REFERENCES claims_officers(id) ON DELETE CASCADE,
	"workload_score" DECIMAL(5,2) NOT NULL,
	"skill_score" DECIMAL(5,2) NOT NULL,
	"total_score" DECIMAL(5,2) NOT NULL,
	"assigned_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assessment_recommendations (
	"id" SERIAL PRIMARY KEY,
	"claim_id" INT UNIQUE NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
	"outcome" VARCHAR(32) NOT NULL,
	"confidence" DECIMAL(5,2) NOT NULL,
	"reasons" TEXT NOT NULL,
	"generated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assessment_recommendations;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS claim_events;
DROP TABLE IF EXISTS claim_documents;
DROP TABLE IF EXISTS claims;
-- +goose StatementEnd

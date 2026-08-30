-- +goose Up
-- +goose StatementBegin
-- Add survey fields to claims table
ALTER TABLE claims 
ADD COLUMN IF NOT EXISTS surveyor_id INT REFERENCES claims_officers(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS survey_status VARCHAR(32),
ADD COLUMN IF NOT EXISTS survey_sla_due_at TIMESTAMPTZ(6),
ADD COLUMN IF NOT EXISTS survey_completed_at TIMESTAMPTZ(6),
ADD COLUMN IF NOT EXISTS survey_report_url TEXT,
ADD COLUMN IF NOT EXISTS survey_photos TEXT[];

-- Add surveyor role and specialty to claims_officers
ALTER TABLE claims_officers
ADD COLUMN IF NOT EXISTS role VARCHAR(64) NOT NULL DEFAULT 'Claims Officer',
ADD COLUMN IF NOT EXISTS specialty VARCHAR(64),
ADD COLUMN IF NOT EXISTS region VARCHAR(64);

-- Update existing officers with roles and specialties
UPDATE claims_officers SET role = 'Claims Officer', specialty = 'Motor', region = 'Central' WHERE id = 1;
UPDATE claims_officers SET role = 'Senior Claims Officer', specialty = 'Motor', region = 'Central' WHERE id = 2;
UPDATE claims_officers SET role = 'Claims Officer', specialty = 'Motor', region = 'North' WHERE id = 3;

-- Add surveyors
INSERT INTO claims_officers (id, name, email, role, specialty, region, current_workload, motor_skill_rating, is_available, created_at, updated_at)
VALUES 
(4, 'Marcus Webb', 'marcus.webb@klemarklemer.com', 'Surveyor', 'Vehicle Inspection', 'Central', 1, 4.70, TRUE, NOW(), NOW()),
(5, 'Priya Sharma', 'priya.sharma@klemarklemer.com', 'Surveyor', 'Property Inspection', 'South', 0, 4.50, TRUE, NOW(), NOW()),
(6, 'James Okafor', 'james.okafor@klemarklemer.com', 'Surveyor', 'Vehicle Inspection', 'North', 2, 4.60, TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Reset sequence
SELECT setval('claims_officers_id_seq', (SELECT MAX(id) FROM claims_officers));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove surveyors
DELETE FROM claims_officers WHERE id IN (4, 5, 6);

-- Reset existing officers
UPDATE claims_officers SET role = 'Claims Officer', specialty = NULL, region = NULL WHERE id IN (1, 2, 3);

-- Remove survey fields from claims
ALTER TABLE claims 
DROP COLUMN IF EXISTS surveyor_id,
DROP COLUMN IF EXISTS survey_status,
DROP COLUMN IF EXISTS survey_sla_due_at,
DROP COLUMN IF EXISTS survey_completed_at,
DROP COLUMN IF EXISTS survey_report_url,
DROP COLUMN IF EXISTS survey_photos;

-- Remove surveyor fields from claims_officers
ALTER TABLE claims_officers
DROP COLUMN IF EXISTS specialty,
DROP COLUMN IF EXISTS region;
-- +goose StatementEnd
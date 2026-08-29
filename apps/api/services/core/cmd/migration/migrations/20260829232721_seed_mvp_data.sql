-- +goose Up
-- +goose StatementBegin
-- 1. Seed Policy
INSERT INTO policies (id, policy_number, policy_holder_name, vehicle_plate, vehicle_model, coverage_type, max_coverage_amount, deductible_amount, effective_date, expiry_date, status, created_at, updated_at)
VALUES (1, 'POL-MOTOR-2026-8819', 'Sarah Jenkins', 'B 1234 KLR', 'Honda CR-V 2023', 'COMPREHENSIVE', 45000.00, 500.00, NOW() - INTERVAL '6 months', NOW() + INTERVAL '6 months', 'ACTIVE', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. Seed 3 Claims Officers with skill vs workload distributions
INSERT INTO claims_officers (id, name, email, role, current_workload, motor_skill_rating, is_available, created_at, updated_at)
VALUES 
(1, 'Alex Rivera', 'alex.rivera@klemarklemer.com', 'Claims Officer', 4, 4.80, TRUE, NOW(), NOW()),
(2, 'David Chen', 'david.chen@klemarklemer.com', 'Senior Claims Officer', 8, 4.90, TRUE, NOW(), NOW()),
(3, 'Elena Rostova', 'elena.rostova@klemarklemer.com', 'Claims Officer', 2, 4.20, TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 3. Seed Incomplete Motor Claim (missing police report)
INSERT INTO claims (id, claim_number, policy_id, stage, document_completeness, survey_required, claim_type, severity, incident_description, estimated_loss, approved_amount, current_officer_id, claim_sla_due_at, stage_sla_due_at, status, created_at, updated_at)
VALUES (
	1, 
	'CLM-2026-0042', 
	1, 
	'DOCUMENT_VERIFICATION', 
	'INCOMPLETE', 
	FALSE, 
	'MOTOR', 
	'MEDIUM', 
	'Front bumper collision with road barrier during heavy rain on highway KM 42.', 
	4200.00, 
	0, 
	NULL, 
	NOW() + INTERVAL '4 hours', 
	NOW() + INTERVAL '25 minutes', 
	'OPEN', 
	NOW() - INTERVAL '15 minutes', 
	NOW()
)
ON CONFLICT (id) DO NOTHING;

-- 4. Seed initial document uploaded (Photo only, missing Police Report)
INSERT INTO claim_documents (id, claim_id, document_type, file_name, file_url, status, extracted_data, uploaded_at)
VALUES (
	1, 
	1, 
	'DAMAGE_PHOTO', 
	'damage_front_bumper_km42.jpg', 
	'https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0042/damage_front_bumper_km42.jpg', 
	'VERIFIED', 
	'{"detected_damage": "Front bumper cracked, grille bent, headlight alignment shifted", "severity": "MEDIUM"}', 
	NOW() - INTERVAL '15 minutes'
)
ON CONFLICT (id) DO NOTHING;

-- 5. Seed initial immutable Claim Events
INSERT INTO claim_events (id, claim_id, actor_name, actor_type, action, previous_stage, new_stage, payload, created_at)
VALUES 
(1, 1, 'Supervisor', 'AGENT', 'CLAIM_INTAKE_INITIALIZED', NULL, 'INTAKE', '{"source": "Customer Submission Portal", "claim_number": "CLM-2026-0042"}', NOW() - INTERVAL '15 minutes'),
(2, 1, 'IntakeAgent', 'AGENT', 'DOCUMENT_VERIFICATION_COMPLETED', 'INTAKE', 'DOCUMENT_VERIFICATION', '{"missing_documents": ["POLICE_REPORT"], "document_completeness": "INCOMPLETE"}', NOW() - INTERVAL '14 minutes')
ON CONFLICT (id) DO NOTHING;

-- Reset sequences
SELECT setval('policies_id_seq', (SELECT MAX(id) FROM policies));
SELECT setval('claims_officers_id_seq', (SELECT MAX(id) FROM claims_officers));
SELECT setval('claims_id_seq', (SELECT MAX(id) FROM claims));
SELECT setval('claim_documents_id_seq', (SELECT MAX(id) FROM claim_documents));
SELECT setval('claim_events_id_seq', (SELECT MAX(id) FROM claim_events));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM claim_events WHERE claim_id = 1;
DELETE FROM claim_documents WHERE claim_id = 1;
DELETE FROM claims WHERE id = 1;
DELETE FROM claims_officers WHERE id IN (1, 2, 3);
DELETE FROM policies WHERE id = 1;
-- +goose StatementEnd

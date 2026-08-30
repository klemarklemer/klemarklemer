-- +goose Up
-- +goose StatementBegin
-- 1. Seed Policy
INSERT INTO policies (id, policy_number, policy_holder_name, vehicle_plate, vehicle_model, coverage_type, max_coverage_amount, deductible_amount, effective_date, expiry_date, status, created_at, updated_at)
VALUES (1, 'POL-MOTOR-2026-8819', 'Sarah Jenkins', 'B 1234 KLR', 'Honda CR-V 2023', 'COMPREHENSIVE', 45000.00, 500.00, NOW() - INTERVAL '6 months', NOW() + INTERVAL '6 months', 'ACTIVE', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. Seed 3 Claims Officers + 3 Surveyors with skill vs workload distributions
INSERT INTO claims_officers (id, name, email, role, specialty, region, current_workload, motor_skill_rating, is_available, created_at, updated_at)
VALUES 
(1, 'Alex Rivera', 'alex.rivera@klemarklemer.com', 'Claims Officer', 'Motor', 'Central', 4, 4.80, TRUE, NOW(), NOW()),
(2, 'David Chen', 'david.chen@klemarklemer.com', 'Senior Claims Officer', 'Motor', 'Central', 8, 4.90, TRUE, NOW(), NOW()),
(3, 'Elena Rostova', 'elena.rostova@klemarklemer.com', 'Claims Officer', 'Motor', 'North', 2, 4.20, TRUE, NOW(), NOW()),
(4, 'Marcus Webb', 'marcus.webb@klemarklemer.com', 'Surveyor', 'Vehicle Inspection', 'Central', 1, 4.70, TRUE, NOW(), NOW()),
(5, 'Priya Sharma', 'priya.sharma@klemarklemer.com', 'Surveyor', 'Property Inspection', 'South', 0, 4.50, TRUE, NOW(), NOW()),
(6, 'James Okafor', 'james.okafor@klemarklemer.com', 'Surveyor', 'Vehicle Inspection', 'North', 2, 4.60, TRUE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 3. Seed Incomplete Motor Claim (missing police report, no survey needed)
INSERT INTO claims (id, claim_number, policy_id, stage, document_completeness, survey_required, surveyor_id, survey_status, survey_sla_due_at, survey_completed_at, survey_report_url, claim_type, severity, incident_description, estimated_loss, approved_amount, current_officer_id, claim_sla_due_at, stage_sla_due_at, status, created_at, updated_at)
VALUES (
	1, 
	'CLM-2026-0042', 
	1, 
	'DOCUMENT_VERIFICATION', 
	'INCOMPLETE', 
	FALSE, 
	NULL, 
	NULL, 
	NULL, 
	NULL, 
	NULL, 
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

-- 4. Seed Survey-Required Motor Claim (with survey in progress)
INSERT INTO claims (id, claim_number, policy_id, stage, document_completeness, survey_required, surveyor_id, survey_status, survey_sla_due_at, survey_completed_at, survey_report_url, claim_type, severity, incident_description, estimated_loss, approved_amount, current_officer_id, claim_sla_due_at, stage_sla_due_at, status, created_at, updated_at)
VALUES (
	2, 
	'CLM-2026-0044', 
	1, 
	'SURVEY', 
	'COMPLETE', 
	TRUE, 
	4, 
	'IN_PROGRESS', 
	NOW() + INTERVAL '2 days', 
	NULL, 
	NULL, 
	'MOTOR', 
	'HIGH', 
	'Multi-vehicle collision on expressway, structural damage suspected', 
	18500.00, 
	0, 
	NULL, 
	NOW() + INTERVAL '7 days', 
	NOW() + INTERVAL '2 days', 
	'OPEN', 
	NOW() - INTERVAL '2 hours', 
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

-- 5. Seed documents for survey-required claim (CLM-2026-0044)
INSERT INTO claim_documents (id, claim_id, document_type, file_name, file_url, status, extracted_data, uploaded_at)
VALUES 
(2, 2, 'DAMAGE_PHOTO', 'collision_scene_01.jpg', 'https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0044/collision_scene_01.jpg', 'VERIFIED', '{"detected_damage": "Front end crumpled, frame rail bent, airbags deployed", "severity": "HIGH"}', NOW() - INTERVAL '2 hours'),
(3, 2, 'POLICE_REPORT', 'police_report_expressway.pdf', 'https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0044/police_report_expressway.pdf', 'VERIFIED', '{"report_id": "PR-2026-1044", "officer": "Sgt. M. Lim", "incident_date": "2026-08-28 14:30", "location": "PIE Expressway KM 12", "vehicles_involved": 3}', NOW() - INTERVAL '1 hour 55 minutes')
ON CONFLICT (id) DO NOTHING;

-- 6. Seed initial immutable Claim Events
INSERT INTO claim_events (id, claim_id, actor_name, actor_type, action, previous_stage, new_stage, payload, created_at)
VALUES 
(1, 1, 'Supervisor', 'AGENT', 'CLAIM_INTAKE_INITIALIZED', NULL, 'INTAKE', '{"source": "Customer Submission Portal", "claim_number": "CLM-2026-0042"}', NOW() - INTERVAL '15 minutes'),
(2, 1, 'IntakeAgent', 'AGENT', 'DOCUMENT_VERIFICATION_COMPLETED', 'INTAKE', 'DOCUMENT_VERIFICATION', '{"missing_documents": ["POLICE_REPORT"], "document_completeness": "INCOMPLETE"}', NOW() - INTERVAL '14 minutes'),
(3, 2, 'Supervisor', 'AGENT', 'CLAIM_INTAKE_INITIALIZED', NULL, 'INTAKE', '{"source": "Customer Submission Portal", "claim_number": "CLM-2026-0044"}', NOW() - INTERVAL '2 hours'),
(4, 2, 'IntakeAgent', 'AGENT', 'DOCUMENT_VERIFICATION_COMPLETED', 'INTAKE', 'DOCUMENT_VERIFICATION', '{"verified_docs": ["DAMAGE_PHOTO", "POLICE_REPORT"], "document_completeness": "COMPLETE"}', NOW() - INTERVAL '1 hour 55 minutes'),
(5, 2, 'IntakeAgent', 'AGENT', 'CLAIM_CLASSIFIED', 'DOCUMENT_VERIFICATION', 'ASSIGNMENT', '{"claim_type": "MOTOR", "severity": "HIGH", "survey_required": true}', NOW() - INTERVAL '1 hour 50 minutes'),
(6, 2, 'AssignmentAgent', 'AGENT', 'SURVEYOR_ASSIGNED', 'ASSIGNMENT', 'SURVEY', '{"surveyor_name": "Marcus Webb", "surveyor_id": 4}', NOW() - INTERVAL '1 hour 45 minutes')
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
DELETE FROM claim_events WHERE claim_id IN (1, 2);
DELETE FROM claim_documents WHERE claim_id IN (1, 2);
DELETE FROM claims WHERE id IN (1, 2);
DELETE FROM claims_officers WHERE id IN (1, 2, 3, 4, 5, 6);
DELETE FROM policies WHERE id = 1;
-- +goose StatementEnd

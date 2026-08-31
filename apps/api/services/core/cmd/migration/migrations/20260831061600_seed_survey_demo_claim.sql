-- +goose Up
-- +goose StatementBegin
-- Survey demo data. Lives here rather than in the MVP seed because the columns
-- it fills are added by 20260830235216, which runs after that seed.
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

INSERT INTO claim_documents (id, claim_id, document_type, file_name, file_url, status, extracted_data, uploaded_at)
VALUES
(2, 2, 'DAMAGE_PHOTO', 'collision_scene_01.jpg', 'https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0044/collision_scene_01.jpg', 'VERIFIED', '{"detected_damage": "Front end crumpled, frame rail bent, airbags deployed", "severity": "HIGH"}', NOW() - INTERVAL '2 hours'),
(3, 2, 'POLICE_REPORT', 'police_report_expressway.pdf', 'https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0044/police_report_expressway.pdf', 'VERIFIED', '{"report_id": "PR-2026-1044", "officer": "Sgt. M. Lim", "incident_date": "2026-08-28 14:30", "location": "PIE Expressway KM 12", "vehicles_involved": 3}', NOW() - INTERVAL '1 hour 55 minutes')
ON CONFLICT (id) DO NOTHING;

INSERT INTO claim_events (id, claim_id, actor_name, actor_type, action, previous_stage, new_stage, payload, created_at)
VALUES
(3, 2, 'Supervisor', 'AGENT', 'CLAIM_INTAKE_INITIALIZED', NULL, 'INTAKE', '{"source": "Customer Submission Portal", "claim_number": "CLM-2026-0044"}', NOW() - INTERVAL '2 hours'),
(4, 2, 'IntakeAgent', 'AGENT', 'DOCUMENT_VERIFICATION_COMPLETED', 'INTAKE', 'DOCUMENT_VERIFICATION', '{"verified_docs": ["DAMAGE_PHOTO", "POLICE_REPORT"], "document_completeness": "COMPLETE"}', NOW() - INTERVAL '1 hour 55 minutes'),
(5, 2, 'IntakeAgent', 'AGENT', 'CLAIM_CLASSIFIED', 'DOCUMENT_VERIFICATION', 'ASSIGNMENT', '{"claim_type": "MOTOR", "severity": "HIGH", "survey_required": true}', NOW() - INTERVAL '1 hour 50 minutes'),
(6, 2, 'AssignmentAgent', 'AGENT', 'SURVEYOR_ASSIGNED', 'ASSIGNMENT', 'SURVEY', '{"surveyor_name": "Marcus Webb", "surveyor_id": 4}', NOW() - INTERVAL '1 hour 45 minutes')
ON CONFLICT (id) DO NOTHING;

SELECT setval('claims_id_seq', (SELECT MAX(id) FROM claims));
SELECT setval('claim_documents_id_seq', (SELECT MAX(id) FROM claim_documents));
SELECT setval('claim_events_id_seq', (SELECT MAX(id) FROM claim_events));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM claim_events WHERE claim_id = 2;
DELETE FROM claim_documents WHERE claim_id = 2;
DELETE FROM claims WHERE id = 2;
-- +goose StatementEnd

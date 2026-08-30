# 20260831012905_add_survey_officers_gemini_taskmaster

## Summary
Added Survey functionality, expanded officer roster with surveyors, and showcased Gemini as the task master orchestrating the assessment flow.

## Changes

### API (Backend)

**New Migration**: `20260830235216_add_survey_and_surveyor.sql`
- Added survey fields to `claims` table:
  - `surveyor_id` (FK to claims_officers)
  - `survey_status` (PENDING, ASSIGNED, IN_PROGRESS, COMPLETED, OVERDUE)
  - `survey_sla_due_at`, `survey_completed_at`, `survey_report_url`, `survey_photos`
- Added `specialty` and `region` columns to `claims_officers` table

**Seed Data Updates** (`20260829232721_seed_mvp_data.sql`):
- Added 3 Claims Officers (Alex Rivera, David Chen, Elena Rostova) with specialties/regions
- Added 3 Surveyors (Marcus Webb - Vehicle Inspection Central, Priya Sharma - Property Inspection South, James Okafor - Vehicle Inspection North)
- Added 2 demo claims:
  - CLM-2026-0042: Incomplete motor claim (no survey required)
  - CLM-2026-0044: Survey-required motor claim with survey in progress

**Domain Models** (`pkg/shared/domain/claim.go`, `officer.go`):
- Added `StageSurvey = "SURVEY"` to claim stages
- Added survey status constants
- Added survey fields to `Claim` struct
- Added `Specialty` and `Region` to `ClaimsOfficer` struct

**Use Cases** (`internal/modules/claim/usecase/`):
- Added `survey.go` with:
  - `AssignSurveyor()` - assigns surveyor, transitions claim to SURVEY stage
  - `UpdateSurveyStatus()` - updates survey status with event logging
  - `CompleteSurvey()` - completes survey, adds survey report/photos as documents, triggers assessment
- Updated `usecase.go` interface with survey methods
- Updated `reset_demo.go` to handle 6 officers and 3 claims

**REST Handlers** (`internal/modules/claim/delivery/resthandler/resthandler.go`):
- Added 3 new endpoints:
  - `POST /v1/claim/:id/survey/assign` - assign surveyor
  - `POST /v1/claim/:id/survey/status` - update survey status
  - `POST /v1/claim/:id/survey/complete` - complete survey and trigger assessment

**Domain Requests** (`internal/modules/claim/domain/request.go`):
- Added `RequestCompleteSurvey` struct with SurveyorID, ReportURL, Photos, Notes

### Web (Frontend)

**Types** (`src/types.ts`):
- Added `Survey` to Stage type
- Added `SurveyInfo` interface
- Added `survey` field to `Claim` interface

**API Client** (`src/api/claims.ts`):
- Added `BackendSurvey` interface
- Added survey fields to `BackendClaim`
- Added `SURVEY` to STAGE_MAP
- Updated `toClaimView()` to map survey data

**Components**:
- New `SurveyPanel.tsx` - displays survey status, surveyor info, report link
- Updated `App.tsx` to include SurveyPanel when claim requires survey
- Updated `RecommendationPanel.tsx` to show "Powered by Gemini" badge showcasing AI task master

**Tests** (`src/App.test.tsx`):
- Added `survey_required: false` to test claim objects

### Scripts
- Updated `demo.sh` with better DSN parsing and `setloss` function

## Demo Flow
1. Run `make up` - starts API on :8000, web on :5173
2. Click "Seed Demo Claim" - creates CLM-2026-0042 (incomplete, no survey) and CLM-2026-0044 (survey in progress)
3. Select CLM-2026-0044 to see Survey panel with Marcus Webb assigned, IN_PROGRESS status
4. Upload police report for CLM-2026-0042 → triggers intake → assignment → assessment
5. Assessment shows "Powered by Gemini" badge - demonstrating Gemini as task master
6. Human approval required for decision (no auto-approve)

## Architecture Notes
- Survey is a new Stage between Assignment and Assessment
- Surveyor assignment uses the same officer pool with role filtering
- Survey completion autonomously triggers Assessment Agent (Loop 3)
- Gemini 3.5 Flash powers the Assessment Agent via the `Assessor` interface
- Deterministic fallback when no credentials configured
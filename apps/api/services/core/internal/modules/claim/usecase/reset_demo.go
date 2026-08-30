package usecase

import (
	"context"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	officerdomain "monorepo/services/core/internal/modules/officer/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) ResetDemo(ctx context.Context) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:ResetDemo")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	var claim shareddomain.Claim
	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		// Reset officers workload
		officers, _ := uc.repoSQL.OfficerRepo().FetchAll(txCtx, &officerdomain.FilterOfficer{})
		for _, off := range officers {
			if off.ID == 1 {
				off.CurrentWorkload = 4
			} else if off.ID == 2 {
				off.CurrentWorkload = 8
			} else if off.ID == 3 {
				off.CurrentWorkload = 2
			}
			_ = uc.repoSQL.OfficerRepo().Save(txCtx, &off)
		}

		// Delete existing demo Claims (and cascade children) before re-seeding
		_ = uc.repoSQL.ClaimRepo().Delete(txCtx, &domain.FilterClaim{ClaimNumber: "CLM-2026-0042"})
		_ = uc.repoSQL.ClaimRepo().Delete(txCtx, &domain.FilterClaim{ClaimNumber: "CLM-2026-0043"})

		claimSLADue := time.Now().Add(4 * time.Hour)
		stageSLADue := time.Now().Add(25 * time.Minute)

		// ---- Claim 1: Normal MOTOR (incomplete, missing police report) ----
		claim = shareddomain.Claim{
			ClaimNumber:          "CLM-2026-0042",
			PolicyID:             1,
			Stage:                shareddomain.StageDocumentVerification,
			DocumentCompleteness: shareddomain.CompletenessIncomplete,
			SurveyRequired:       false,
			ClaimType:            "MOTOR",
			Severity:             "MEDIUM",
			IncidentDescription:  "Front bumper collision with road barrier during heavy rain on highway KM 42.",
			EstimatedLoss:        4200.00,
			ApprovedAmount:       0,
			CurrentOfficerID:     nil,
			ClaimSLADueAt:        &claimSLADue,
			StageSLADueAt:        &stageSLADue,
			Status:               "OPEN",
			CreatedAt:            time.Now().Add(-15 * time.Minute),
			UpdatedAt:            time.Now(),
		}

		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		doc := shareddomain.ClaimDocument{
			ClaimID:       claim.ID,
			DocumentType:  "DAMAGE_PHOTO",
			FileName:      "damage_front_bumper_km42.jpg",
			FileURL:       "https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0042/damage_front_bumper_km42.jpg",
			Status:        "VERIFIED",
			ExtractedData: `{"detected_damage": "Front bumper cracked, grille bent, headlight alignment shifted", "severity": "MEDIUM"}`,
			UploadedAt:    time.Now().Add(-15 * time.Minute),
		}
		if err := uc.repoSQL.ClaimRepo().AddDocument(txCtx, &doc); err != nil {
			return err
		}

		event1 := shareddomain.ClaimEvent{
			ClaimID:   claim.ID,
			ActorName: "Supervisor",
			ActorType: shareddomain.ActorAgent,
			Action:    "CLAIM_INTAKE_INITIALIZED",
			NewStage:  shareddomain.StageIntake,
			Payload:   `{"source": "Customer Submission Portal", "claim_number": "CLM-2026-0042"}`,
			CreatedAt: time.Now().Add(-15 * time.Minute),
		}
		_ = uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event1)

		event2 := shareddomain.ClaimEvent{
			ClaimID:       claim.ID,
			ActorName:     "IntakeAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "DOCUMENT_VERIFICATION_COMPLETED",
			PreviousStage: shareddomain.StageIntake,
			NewStage:      shareddomain.StageDocumentVerification,
			Payload:       `{"missing_documents": ["POLICE_REPORT"], "document_completeness": "INCOMPLETE"}`,
			CreatedAt:     time.Now().Add(-14 * time.Minute),
		}
		if err := uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event2); err != nil {
			return err
		}

		// ---- Claim 2: Fraud-flagged MOTOR (policy holder mismatch, mandatory review) ----
		fraudSignal := "Policy holder name on claim (J. Tan) differs from policy (Sarah Jenkins); duplicate claim pattern detected; vehicle plate B 1234 KLR not registered to claimant"
		fraudStageSLADue := time.Now().Add(10 * time.Minute) // at risk
		fraudClaim := shareddomain.Claim{
			ClaimNumber:          "CLM-2026-0043",
			PolicyID:             1,
			Stage:                shareddomain.StageDecision,
			DocumentCompleteness: shareddomain.CompletenessComplete,
			SurveyRequired:       false,
			ClaimType:            "MOTOR",
			Severity:             "HIGH",
			IncidentDescription:  "Alleged theft of vehicle from residential driveway.",
			EstimatedLoss:        45000.00,
			ApprovedAmount:       0,
			CurrentOfficerID:     nil,
			FraudSignal:          &fraudSignal,
			ClaimSLADueAt:        &claimSLADue,
			StageSLADueAt:        &fraudStageSLADue,
			Status:               "OPEN",
			CreatedAt:            time.Now().Add(-45 * time.Minute),
			UpdatedAt:            time.Now(),
		}

		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &fraudClaim); err != nil {
			return err
		}

		// Inconsistent claim form document
		doc2 := shareddomain.ClaimDocument{
			ClaimID:       fraudClaim.ID,
			DocumentType:  "CLAIM_FORM",
			FileName:      "claim_form_theft_report.pdf",
			FileURL:       "https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0043/claim_form_theft_report.pdf",
			Status:        "VERIFIED",
			ExtractedData: `{"claimant_name": "J. Tan", "claimant_id": "S1234567", "vehicle_plate": "B 1234 KLR", "incident_type": "THEFT", "reported_date": "2026-08-28"}`,
			UploadedAt:    time.Now().Add(-40 * time.Minute),
		}
		if err := uc.repoSQL.ClaimRepo().AddDocument(txCtx, &doc2); err != nil {
			return err
		}

		// Police report for fraud claim
		doc3 := shareddomain.ClaimDocument{
			ClaimID:       fraudClaim.ID,
			DocumentType:  "POLICE_REPORT",
			FileName:      "police_report_theft.pdf",
			FileURL:       "https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0043/police_report_theft.pdf",
			Status:        "VERIFIED",
			ExtractedData: `{"report_id": "PR-2026-9999", "officer": "Insp. A. Kumar", "incident_date": "2026-08-28 02:00", "location": "123 Jalan Melati", "status": "OPEN_INVESTIGATION"}`,
			UploadedAt:    time.Now().Add(-35 * time.Minute),
		}
		if err := uc.repoSQL.ClaimRepo().AddDocument(txCtx, &doc3); err != nil {
			return err
		}

		// Mandatory human review recommendation (not auto-approve)
		rec := shareddomain.AssessmentRecommendation{
			ClaimID:     fraudClaim.ID,
			Outcome:     shareddomain.OutcomeManualReview,
			Confidence:  0.87,
			Reasons:     "Claimant name (J. Tan) does not match policy holder (Sarah Jenkins). Vehicle plate B 1234 KLR registered to policy holder, not claimant. Prior theft claim on same vehicle 6 months ago (CLM-2025-1102). High estimated loss ($45,000) at policy maximum. Fraud signals require investigation before any decision.",
			GeneratedAt: time.Now(),
		}
		if err := uc.repoSQL.ClaimRepo().SaveRecommendation(txCtx, &rec); err != nil {
			return err
		}

		fraudEvent := shareddomain.ClaimEvent{
			ClaimID:       fraudClaim.ID,
			ActorName:     "AssessmentAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "FRAUD_SIGNAL_DETECTED",
			PreviousStage: shareddomain.StageAssessment,
			NewStage:      shareddomain.StageDecision,
			Payload:       `{"fraud_signal": "Policy holder mismatch; duplicate vehicle claim history", "recommendation": "MANUAL_REVIEW", "confidence": 0.87}`,
			CreatedAt:     time.Now().Add(-5 * time.Minute),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &fraudEvent)
	})

	if err != nil {
		return res, err
	}

	// Return the normal claim (CLM-2026-0042) as the seeded workspace
	return uc.GetDetailClaim(ctx, claim.ID)
}
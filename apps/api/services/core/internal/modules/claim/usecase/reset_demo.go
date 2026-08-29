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

		// Delete existing documents, events, assignments, recommendations for claim 1
		_ = uc.repoSQL.ClaimRepo().Delete(txCtx, &domain.FilterClaim{ID: func() *int { id := 1; return &id }()})

		claimSLADue := time.Now().Add(4 * time.Hour)
		stageSLADue := time.Now().Add(25 * time.Minute)

		claim := shareddomain.Claim{
			ID:                   1,
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
			ClaimID:       1,
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
			ClaimID:   1,
			ActorName: "Supervisor",
			ActorType: shareddomain.ActorAgent,
			Action:    "CLAIM_INTAKE_INITIALIZED",
			NewStage:  shareddomain.StageIntake,
			Payload:   `{"source": "Customer Submission Portal", "claim_number": "CLM-2026-0042"}`,
			CreatedAt: time.Now().Add(-15 * time.Minute),
		}
		_ = uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event1)

		event2 := shareddomain.ClaimEvent{
			ClaimID:       1,
			ActorName:     "IntakeAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "DOCUMENT_VERIFICATION_COMPLETED",
			PreviousStage: shareddomain.StageIntake,
			NewStage:      shareddomain.StageDocumentVerification,
			Payload:       `{"missing_documents": ["POLICE_REPORT"], "document_completeness": "INCOMPLETE"}`,
			CreatedAt:     time.Now().Add(-14 * time.Minute),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event2)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, 1)
}

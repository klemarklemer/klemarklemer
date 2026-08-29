package usecase

import (
	"context"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) EvaluateIntake(ctx context.Context, claimID int) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:EvaluateIntake")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	docs, err := uc.repoSQL.ClaimRepo().GetDocumentsByClaimID(ctx, claimID)
	if err != nil {
		return res, err
	}

	hasPoliceReport := false
	hasDamagePhoto := false
	for _, doc := range docs {
		if doc.DocumentType == "POLICE_REPORT" {
			hasPoliceReport = true
		}
		if doc.DocumentType == "DAMAGE_PHOTO" {
			hasDamagePhoto = true
		}
	}

	if hasPoliceReport && hasDamagePhoto {
		// Complete documents
		claim.DocumentCompleteness = shareddomain.CompletenessComplete
		claim.SurveyRequired = false
		claim.Severity = "MEDIUM"
		claim.Stage = shareddomain.StageAssignment
		newStageSLA := time.Now().Add(20 * time.Minute)
		claim.StageSLADueAt = &newStageSLA

		err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
			if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
				return err
			}

			event1 := shareddomain.ClaimEvent{
				ClaimID:       claimID,
				ActorName:     "IntakeAgent",
				ActorType:     shareddomain.ActorAgent,
				Action:        "DOCUMENT_VERIFICATION_COMPLETED",
				PreviousStage: shareddomain.StageDocumentVerification,
				NewStage:      shareddomain.StageAssignment,
				Payload:       `{"document_completeness":"COMPLETE","verified_docs":["POLICE_REPORT","DAMAGE_PHOTO"]}`,
			}
			if err := uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event1); err != nil {
				return err
			}

			event2 := shareddomain.ClaimEvent{
				ClaimID:   claimID,
				ActorName: "IntakeAgent",
				ActorType: shareddomain.ActorAgent,
				Action:    "CLAIM_CLASSIFIED",
				Payload:   `{"claim_type":"MOTOR","severity":"MEDIUM","survey_required":false}`,
			}
			return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event2)
		})

		if err != nil {
			return res, err
		}

		// Autonomous progression: Trigger Assignment Agent
		return uc.RunAssignment(ctx, claimID)
	}

	// Still incomplete
	claim.DocumentCompleteness = shareddomain.CompletenessIncomplete
	claim.Stage = shareddomain.StageDocumentVerification

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "IntakeAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "DOCUMENT_COMPLETENESS_CHECKED",
			PreviousStage: shareddomain.StageIntake,
			NewStage:      shareddomain.StageDocumentVerification,
			Payload:       `{"document_completeness":"INCOMPLETE","missing_documents":["POLICE_REPORT"]}`,
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

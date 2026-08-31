package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	officerdomain "monorepo/services/core/internal/modules/officer/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"
	"monorepo/services/core/pkg/shared/gemini"

	"github.com/golangid/candi/tracer"
)

// EvaluateIntake is Loop 1. It reads the claim and the documents on file, has the
// Intake Agent classify severity and whether a surveyor is needed, and records
// which required documents are genuinely absent.
//
// A complete claim chains onward: to Loop 2b when the agent asks for an
// inspection, otherwise straight to officer assignment.
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

	classification, err := uc.classifier.Classify(ctx, buildClassificationInput(claim, docs))
	if err != nil {
		return res, fmt.Errorf("classify claim %d: %w", claimID, err)
	}

	claim.Severity = classification.Severity
	claim.SurveyRequired = classification.SurveyRequired

	if len(classification.MissingDocuments) > 0 {
		return uc.recordIncompleteIntake(ctx, claim, classification)
	}
	return uc.recordCompleteIntake(ctx, claim, classification, documentTypes(docs))
}

func (uc *claimUsecaseImpl) recordCompleteIntake(
	ctx context.Context, claim shareddomain.Claim, c gemini.ClassificationResult, present []string,
) (res domain.ResponseClaim, err error) {
	claimID := claim.ID
	claim.DocumentCompleteness = shareddomain.CompletenessComplete
	claim.Stage = shareddomain.StageAssignment
	newStageSLA := time.Now().Add(20 * time.Minute)
	claim.StageSLADueAt = &newStageSLA

	verified, err := json.Marshal(map[string]any{
		"document_completeness": shareddomain.CompletenessComplete,
		"verified_docs":         present,
	})
	if err != nil {
		return res, fmt.Errorf("encode intake verification for claim %d: %w", claimID, err)
	}

	classified, err := json.Marshal(map[string]any{
		"claim_type":      claim.ClaimType,
		"severity":        c.Severity,
		"survey_required": c.SurveyRequired,
		"reasons":         c.Reasons,
		"engine":          c.Source,
	})
	if err != nil {
		return res, fmt.Errorf("encode intake classification for claim %d: %w", claimID, err)
	}

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		if err := uc.repoSQL.ClaimRepo().AddEvent(txCtx, &shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "IntakeAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "DOCUMENT_VERIFICATION_COMPLETED",
			PreviousStage: shareddomain.StageDocumentVerification,
			NewStage:      shareddomain.StageAssignment,
			Payload:       string(verified),
		}); err != nil {
			return err
		}

		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &shareddomain.ClaimEvent{
			ClaimID:   claimID,
			ActorName: "IntakeAgent",
			ActorType: shareddomain.ActorAgent,
			Action:    "CLAIM_CLASSIFIED",
			Payload:   string(classified),
		})
	})
	if err != nil {
		return res, err
	}

	if c.SurveyRequired {
		return uc.routeToSurvey(ctx, claimID, claim.ClaimType)
	}
	return uc.RunAssignment(ctx, claimID)
}

func (uc *claimUsecaseImpl) recordIncompleteIntake(
	ctx context.Context, claim shareddomain.Claim, c gemini.ClassificationResult,
) (res domain.ResponseClaim, err error) {
	claimID := claim.ID
	claim.DocumentCompleteness = shareddomain.CompletenessIncomplete
	claim.Stage = shareddomain.StageDocumentVerification

	payload, err := json.Marshal(map[string]any{
		"document_completeness": shareddomain.CompletenessIncomplete,
		"missing_documents":     c.MissingDocuments,
		"severity":              c.Severity,
		"reasons":               c.Reasons,
		"engine":                c.Source,
	})
	if err != nil {
		return res, fmt.Errorf("encode intake completeness for claim %d: %w", claimID, err)
	}

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "IntakeAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "DOCUMENT_COMPLETENESS_CHECKED",
			PreviousStage: shareddomain.StageIntake,
			NewStage:      shareddomain.StageDocumentVerification,
			Payload:       string(payload),
		})
	})
	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

// routeToSurvey hands a claim the Intake Agent flagged to an available surveyor.
// When none is free the claim still advances to officer assignment rather than
// stalling - a surveyor who is merely busy must not strand the workflow.
func (uc *claimUsecaseImpl) routeToSurvey(ctx context.Context, claimID int, claimType string) (res domain.ResponseClaim, err error) {
	officers, err := uc.repoSQL.OfficerRepo().FetchAll(ctx, &officerdomain.FilterOfficer{})
	if err != nil {
		return res, fmt.Errorf("load surveyors: %w", err)
	}

	surveyor := selectSurveyor(officers, claimType)
	if surveyor == nil {
		return uc.RunAssignment(ctx, claimID)
	}
	return uc.AssignSurveyor(ctx, claimID, surveyor.ID)
}

func buildClassificationInput(claim shareddomain.Claim, docs []shareddomain.ClaimDocument) gemini.ClassificationInput {
	in := gemini.ClassificationInput{
		ClaimNumber:         claim.ClaimNumber,
		ClaimType:           claim.ClaimType,
		IncidentDescription: claim.IncidentDescription,
		EstimatedLoss:       claim.EstimatedLoss,
		DocumentTypes:       documentTypes(docs),
	}
	if claim.Policy != nil {
		in.MaxCoverageAmount = claim.Policy.MaxCoverageAmount
	}
	return in
}

func documentTypes(docs []shareddomain.ClaimDocument) []string {
	types := make([]string, 0, len(docs))
	for _, doc := range docs {
		types = append(types, doc.DocumentType)
	}
	return types
}

package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"
	"monorepo/services/core/pkg/shared/gemini"

	"github.com/golangid/candi/tracer"
)

// decisionStageSLA is how long the human decision gate has to bind a
// recommendation before the claim breaches its stage SLA.
const decisionStageSLA = 10 * time.Minute

// RunAssessment is Loop 3. It asks the Assessment Agent to weigh the policy's
// coverage against the estimated loss and the documents on file, records the
// recommendation, and hands the claim to the human decision gate.
//
// The recommendation is advisory. Nothing here settles or closes a claim - that
// is SubmitHumanApproval's job - so an APPROVE from this loop is a proposal a
// claims officer still has to bind.
func (uc *claimUsecaseImpl) RunAssessment(ctx context.Context, claimID int) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:RunAssessment")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	assessment, err := uc.assessor.Assess(ctx, buildAssessmentInput(claim))
	if err != nil {
		return res, fmt.Errorf("assess claim %d: %w", claimID, err)
	}

	rec := shareddomain.AssessmentRecommendation{
		ClaimID:     claimID,
		Outcome:     assessment.Outcome,
		Confidence:  assessment.Confidence,
		Reasons:     assessment.Reasons,
		GeneratedAt: time.Now(),
	}

	claim.Stage = shareddomain.StageDecision
	newStageSLA := time.Now().Add(decisionStageSLA)
	claim.StageSLADueAt = &newStageSLA

	payload, err := json.Marshal(map[string]any{
		"outcome":    assessment.Outcome,
		"confidence": assessment.Confidence,
		"reasons":    assessment.Reasons,
		"engine":     assessment.Source,
	})
	if err != nil {
		return res, fmt.Errorf("encode assessment event for claim %d: %w", claimID, err)
	}

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().SaveRecommendation(txCtx, &rec); err != nil {
			return err
		}

		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "AssessmentAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "RECOMMENDATION_GENERATED",
			PreviousStage: shareddomain.StageAssessment,
			NewStage:      shareddomain.StageDecision,
			Payload:       string(payload),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

// buildAssessmentInput flattens the claim record into the evidence the agent is
// allowed to reason over. A claim with no policy attached still produces a valid
// input; the assessor is responsible for treating that as unestablished coverage
// rather than assuming a default.
func buildAssessmentInput(claim shareddomain.Claim) gemini.AssessmentInput {
	in := gemini.AssessmentInput{
		ClaimNumber:         claim.ClaimNumber,
		ClaimType:           claim.ClaimType,
		Severity:            claim.Severity,
		IncidentDescription: claim.IncidentDescription,
		EstimatedLoss:       claim.EstimatedLoss,
	}

	if claim.Policy != nil {
		in.PolicyNumber = claim.Policy.PolicyNumber
		in.CoverageType = claim.Policy.CoverageType
		in.PolicyStatus = claim.Policy.Status
		in.MaxCoverageAmount = claim.Policy.MaxCoverageAmount
		in.DeductibleAmount = claim.Policy.DeductibleAmount
		in.PolicyExpiry = claim.Policy.ExpiryDate
	}

	for _, doc := range claim.Documents {
		in.DocumentTypes = append(in.DocumentTypes, doc.DocumentType)
	}

	return in
}

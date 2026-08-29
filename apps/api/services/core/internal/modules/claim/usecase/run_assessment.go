package usecase

import (
	"context"
	"fmt"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) RunAssessment(ctx context.Context, claimID int) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:RunAssessment")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	reasons := "Damage matches single-vehicle road barrier collision documented in Police Report PR-2026-9912. Policy #POL-MOTOR-2026-8819 is in-force with Comprehensive motor coverage up to $45,000.00. Estimated loss of $4,200.00 is fully substantiated by photo telemetry and within allowable limits after $500.00 deductible."

	rec := shareddomain.AssessmentRecommendation{
		ClaimID:     claimID,
		Outcome:     shareddomain.OutcomeApprove,
		Confidence:  0.94,
		Reasons:     reasons,
		GeneratedAt: time.Now(),
	}

	claim.Stage = shareddomain.StageDecision
	newStageSLA := time.Now().Add(10 * time.Minute)
	claim.StageSLADueAt = &newStageSLA

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
			Payload: fmt.Sprintf(
				`{"outcome":"APPROVE","confidence":0.94,"reasons":"%s"}`,
				rec.Reasons,
			),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

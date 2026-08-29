package usecase

import (
	"context"
	"fmt"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) CreateClaim(ctx context.Context, data *domain.RequestCreateClaim) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:CreateClaim")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	claimNumber := fmt.Sprintf("CLM-2026-%04d", time.Now().Unix()%10000)
	now := time.Now()
	claimSLADue := now.Add(4 * time.Hour)
	stageSLADue := now.Add(30 * time.Minute)

	claim := shareddomain.Claim{
		ClaimNumber:          claimNumber,
		PolicyID:             data.PolicyID,
		Stage:                shareddomain.StageIntake,
		DocumentCompleteness: shareddomain.CompletenessIncomplete,
		SurveyRequired:       false,
		ClaimType:            data.ClaimType,
		Severity:             "MEDIUM",
		IncidentDescription:  data.IncidentDescription,
		EstimatedLoss:        data.EstimatedLoss,
		Status:               "OPEN",
		ClaimSLADueAt:        &claimSLADue,
		StageSLADueAt:        &stageSLADue,
	}

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:   claim.ID,
			ActorName: "Supervisor",
			ActorType: shareddomain.ActorAgent,
			Action:    "CLAIM_CREATED",
			NewStage:  shareddomain.StageIntake,
			Payload:   fmt.Sprintf(`{"claim_number":"%s","policy_id":%d}`, claim.ClaimNumber, claim.PolicyID),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claim.ID)
}

func (uc *claimUsecaseImpl) UpdateClaim(ctx context.Context, data *domain.RequestClaim) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:UpdateClaim")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	claim := data.Deserialize()
	return uc.repoSQL.ClaimRepo().Save(ctx, &claim)
}

func (uc *claimUsecaseImpl) DeleteClaim(ctx context.Context, id int) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:DeleteClaim")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	return uc.repoSQL.ClaimRepo().Delete(ctx, &domain.FilterClaim{ID: &id})
}

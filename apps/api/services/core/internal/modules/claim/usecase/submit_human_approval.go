package usecase

import (
	"context"
	"fmt"

	"monorepo/services/core/internal/modules/claim/domain"
	officerdomain "monorepo/services/core/internal/modules/officer/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) SubmitHumanApproval(ctx context.Context, claimID int, data *domain.RequestHumanApproval) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:SubmitHumanApproval")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	officerName := "Claims Officer"
	if data.OfficerID > 0 {
		offFilter := officerdomain.FilterOfficer{ID: &data.OfficerID}
		if off, err := uc.repoSQL.OfficerRepo().Find(ctx, &offFilter); err == nil {
			officerName = off.Name
		}
	} else if claim.CurrentOfficer != nil {
		officerName = claim.CurrentOfficer.Name
	}

	if data.Action == "APPROVE" {
		deductible := 500.00
		if claim.Policy != nil && claim.Policy.DeductibleAmount > 0 {
			deductible = claim.Policy.DeductibleAmount
		}
		claim.ApprovedAmount = claim.EstimatedLoss - deductible
		if claim.ApprovedAmount < 0 {
			claim.ApprovedAmount = 0
		}
		claim.Stage = shareddomain.StageClosed
		claim.Status = "CLOSED"

		err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
			if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
				return err
			}

			// Decrement officer workload upon claim closure
			if claim.CurrentOfficerID != nil {
				offID := *claim.CurrentOfficerID
				offFilter := officerdomain.FilterOfficer{ID: &offID}
				if off, err := uc.repoSQL.OfficerRepo().Find(txCtx, &offFilter); err == nil && off.CurrentWorkload > 0 {
					off.CurrentWorkload--
					_ = uc.repoSQL.OfficerRepo().Save(txCtx, &off)
				}
			}

			event1 := shareddomain.ClaimEvent{
				ClaimID:       claimID,
				ActorName:     officerName,
				ActorType:     shareddomain.ActorOfficer,
				Action:        "HUMAN_APPROVAL_RECORDED",
				PreviousStage: shareddomain.StageDecision,
				NewStage:      shareddomain.StageClosed,
				Payload: fmt.Sprintf(
					`{"action":"APPROVE","officer_name":"%s","notes":"%s","approved_amount":%.2f}`,
					officerName, data.Notes, claim.ApprovedAmount,
				),
			}
			if err := uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event1); err != nil {
				return err
			}

			event2 := shareddomain.ClaimEvent{
				ClaimID:   claimID,
				ActorName: "Supervisor",
				ActorType: shareddomain.ActorAgent,
				Action:    "DECISION_ISSUED",
				Payload: fmt.Sprintf(
					`{"binding_outcome":"APPROVE","settlement_amount":%.2f,"status":"CLOSED","generated_report":"claim_assessment_report_CLM-2026-0042.pdf"}`,
					claim.ApprovedAmount,
				),
			}
			return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event2)
		})
	} else {
		// Rejection / Override
		claim.Stage = shareddomain.StageAssessment

		err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
			if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
				return err
			}

			event := shareddomain.ClaimEvent{
				ClaimID:       claimID,
				ActorName:     officerName,
				ActorType:     shareddomain.ActorOfficer,
				Action:        "HUMAN_APPROVAL_REJECTED",
				PreviousStage: shareddomain.StageDecision,
				NewStage:      shareddomain.StageAssessment,
				Payload: fmt.Sprintf(
					`{"action":"REJECT","officer_name":"%s","reasons":"%s"}`,
					officerName, data.Notes,
				),
			}
			return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
		})
	}

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

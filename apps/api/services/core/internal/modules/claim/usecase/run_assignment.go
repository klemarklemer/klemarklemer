package usecase

import (
	"context"
	"fmt"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	officerdomain "monorepo/services/core/internal/modules/officer/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) RunAssignment(ctx context.Context, claimID int) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:RunAssignment")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	officers, err := uc.repoSQL.OfficerRepo().FetchAll(ctx, &officerdomain.FilterOfficer{})
	if err != nil {
		return res, fmt.Errorf("load claims officers: %w", err)
	}
	if len(officers) == 0 {
		return res, fmt.Errorf("no claims officers exist to assign claim %d", claimID)
	}

	var bestOfficer *shareddomain.Officer
	var highestScore float64 = -1
	var bestWorkloadScore, bestSkillScore float64

	for i := range officers {
		off := &officers[i]
		if !off.IsAvailable {
			continue
		}

		workloadCap := 10.0
		workloadDiff := workloadCap - float64(off.CurrentWorkload)
		if workloadDiff < 0 {
			workloadDiff = 0
		}
		workloadScore := workloadDiff * 0.5
		skillScore := (off.MotorSkillRating / 5.0) * 10.0 * 0.5
		totalScore := workloadScore + skillScore

		if totalScore > highestScore {
			highestScore = totalScore
			bestOfficer = off
			bestWorkloadScore = workloadScore
			bestSkillScore = skillScore
		}
	}

	if bestOfficer == nil {
		bestOfficer = &officers[0]
		highestScore = 5.0
	}

	assignment := shareddomain.Assignment{
		ClaimID:       claimID,
		OfficerID:     bestOfficer.ID,
		WorkloadScore: bestWorkloadScore,
		SkillScore:    bestSkillScore,
		TotalScore:    highestScore,
		AssignedAt:    time.Now(),
	}

	claim.CurrentOfficerID = &bestOfficer.ID
	claim.Stage = shareddomain.StageAssessment
	newStageSLA := time.Now().Add(15 * time.Minute)
	claim.StageSLADueAt = &newStageSLA

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().SaveAssignment(txCtx, &assignment); err != nil {
			return err
		}

		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		// Update officer workload
		bestOfficer.CurrentWorkload++
		if err := uc.repoSQL.OfficerRepo().Save(txCtx, bestOfficer); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "AssignmentAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "OFFICER_ASSIGNED",
			PreviousStage: shareddomain.StageAssignment,
			NewStage:      shareddomain.StageAssessment,
			Payload: fmt.Sprintf(
				`{"officer_name":"%s","officer_id":%d,"workload_score":%.2f,"skill_score":%.2f,"total_score":%.2f}`,
				bestOfficer.Name, bestOfficer.ID, bestWorkloadScore, bestSkillScore, highestScore,
			),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	// Autonomous progression: Trigger Assessment Agent
	return uc.RunAssessment(ctx, claimID)
}

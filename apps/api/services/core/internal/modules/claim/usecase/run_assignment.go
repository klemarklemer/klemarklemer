package usecase

import (
	"context"
	"fmt"
	"strings"
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

	bestOfficer, bestWorkloadScore, bestSkillScore, highestScore := selectClaimOwner(officers)
	if bestOfficer == nil {
		return res, fmt.Errorf("no available claims officer to assign claim %d", claimID)
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

// selectClaimOwner picks the available claims officer with the best combined
// workload and skill score. Surveyors sit in the same table but inspect damage
// rather than owning claims, so they are never candidates - without this filter
// they outscore every officer and Loop 2 hands motor claims to an inspector.
// Returns a nil officer when nobody is eligible.
func selectClaimOwner(officers []shareddomain.Officer) (best *shareddomain.Officer, workload, skill, total float64) {
	total = -1
	for i := range officers {
		off := &officers[i]
		if !off.IsAvailable || off.IsSurveyor() {
			continue
		}

		workloadCap := 10.0
		workloadDiff := workloadCap - float64(off.CurrentWorkload)
		if workloadDiff < 0 {
			workloadDiff = 0
		}
		workloadScore := workloadDiff * 0.5
		skillScore := (off.MotorSkillRating / 5.0) * 10.0 * 0.5

		if workloadScore+skillScore > total {
			best, workload, skill, total = off, workloadScore, skillScore, workloadScore+skillScore
		}
	}

	if best == nil {
		return nil, 0, 0, 0
	}
	return best, workload, skill, total
}

// surveyorSpecialty names the inspection specialty that fits a claim type. A
// property inspector sent to a car crash is the same class of mistake as a
// surveyor owning a claim, so the fit is chosen before the workload.
var surveyorSpecialty = map[string]string{
	"MOTOR": "Vehicle Inspection",
}

// selectSurveyor picks the available surveyor best suited to the claim: the
// least-loaded one whose specialty fits, falling back to the least-loaded
// surveyor of any specialty. Skill rating scores claim ownership rather than
// inspection, so it is deliberately not used here. Returns nil when none is free.
func selectSurveyor(officers []shareddomain.Officer, claimType string) *shareddomain.Officer {
	preferred := surveyorSpecialty[strings.ToUpper(strings.TrimSpace(claimType))]

	var bestFit, bestAny *shareddomain.Officer
	for i := range officers {
		off := &officers[i]
		if !off.IsAvailable || !off.IsSurveyor() {
			continue
		}
		if bestAny == nil || off.CurrentWorkload < bestAny.CurrentWorkload {
			bestAny = off
		}
		if preferred != "" && off.Specialty != nil && strings.EqualFold(*off.Specialty, preferred) {
			if bestFit == nil || off.CurrentWorkload < bestFit.CurrentWorkload {
				bestFit = off
			}
		}
	}

	if bestFit != nil {
		return bestFit
	}
	return bestAny
}

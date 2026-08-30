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

// AssignSurveyor assigns a surveyor to a claim
func (uc *claimUsecaseImpl) AssignSurveyor(ctx context.Context, claimID int, surveyorID int) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:AssignSurveyor")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	// Verify surveyor exists and is available
	surveyorFilter := officerdomain.FilterOfficer{ID: &surveyorID}
	surveyor, err := uc.repoSQL.OfficerRepo().Find(ctx, &surveyorFilter)
	if err != nil {
		return res, fmt.Errorf("surveyor not found: %w", err)
	}

	if !surveyor.IsAvailable {
		return res, fmt.Errorf("surveyor %s is not available", surveyor.Name)
	}

	// Verify surveyor is actually a surveyor
	if surveyor.Role != "Surveyor" {
		return res, fmt.Errorf("officer %s is not a surveyor", surveyor.Name)
	}

	claim.SurveyorID = &surveyorID
	claim.SurveyStatus = shareddomain.SurveyStatusAssigned
	claim.Stage = shareddomain.StageSurvey
	surveySLADue := time.Now().Add(48 * time.Hour)
	claim.SurveySLADueAt = &surveySLADue

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "AssignmentAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "SURVEYOR_ASSIGNED",
			PreviousStage: shareddomain.StageAssignment,
			NewStage:      shareddomain.StageSurvey,
			Payload:       fmt.Sprintf(`{"surveyor_name":"%s","surveyor_id":%d,"survey_sla_due":"%s"}`, surveyor.Name, surveyorID, surveySLADue.Format(time.RFC3339)),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

// UpdateSurveyStatus updates the survey status
func (uc *claimUsecaseImpl) UpdateSurveyStatus(ctx context.Context, claimID int, status string) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:UpdateSurveyStatus")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	validStatuses := map[string]bool{
		shareddomain.SurveyStatusPending:    true,
		shareddomain.SurveyStatusAssigned:   true,
		shareddomain.SurveyStatusInProgress: true,
		shareddomain.SurveyStatusCompleted:  true,
		shareddomain.SurveyStatusOverdue:    true,
	}

	if !validStatuses[status] {
		return res, fmt.Errorf("invalid survey status: %s", status)
	}

	claim.SurveyStatus = status

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:   claimID,
			ActorName: "SurveyAgent",
			ActorType: shareddomain.ActorAgent,
			Action:    "SURVEY_STATUS_UPDATED",
			Payload:   fmt.Sprintf(`{"survey_status":"%s"}`, status),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	return uc.GetDetailClaim(ctx, claimID)
}

// CompleteSurvey completes the survey and moves to assessment
func (uc *claimUsecaseImpl) CompleteSurvey(ctx context.Context, claimID int, data *domain.RequestCompleteSurvey) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:CompleteSurvey")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &claimID}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return res, err
	}

	if claim.SurveyRequired == false {
		return res, fmt.Errorf("survey not required for this claim")
	}

	if claim.SurveyorID == nil || *claim.SurveyorID != data.SurveyorID {
		return res, fmt.Errorf("surveyor mismatch")
	}

	now := time.Now()
	claim.SurveyStatus = shareddomain.SurveyStatusCompleted
	claim.SurveyCompletedAt = &now
	claim.SurveyReportURL = data.ReportURL
	claim.SurveyPhotos = data.Photos
	claim.Stage = shareddomain.StageAssessment

	stageSLADue := time.Now().Add(15 * time.Minute)
	claim.StageSLADueAt = &stageSLADue

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().Save(txCtx, &claim); err != nil {
			return err
		}

		// Add survey report as document
		doc := shareddomain.ClaimDocument{
			ClaimID:       claimID,
			DocumentType:  "SURVEY_REPORT",
			FileName:      "survey_report.pdf",
			FileURL:       data.ReportURL,
			Status:        "VERIFIED",
			ExtractedData: data.Notes,
			UploadedAt:    now,
		}
		if err := uc.repoSQL.ClaimRepo().AddDocument(txCtx, &doc); err != nil {
			return err
		}

		// Add survey photos as documents
		for i, photoURL := range data.Photos {
			photoDoc := shareddomain.ClaimDocument{
				ClaimID:       claimID,
				DocumentType:  "SURVEY_PHOTO",
				FileName:      fmt.Sprintf("survey_photo_%d.jpg", i+1),
				FileURL:       photoURL,
				Status:        "VERIFIED",
				UploadedAt:    now,
			}
			if err := uc.repoSQL.ClaimRepo().AddDocument(txCtx, &photoDoc); err != nil {
				return err
			}
		}

		event := shareddomain.ClaimEvent{
			ClaimID:       claimID,
			ActorName:     "SurveyAgent",
			ActorType:     shareddomain.ActorAgent,
			Action:        "SURVEY_COMPLETED",
			PreviousStage: shareddomain.StageSurvey,
			NewStage:      shareddomain.StageAssessment,
			Payload:       fmt.Sprintf(`{"surveyor_id":%d,"survey_report_url":"%s","photos_count":%d}`, data.SurveyorID, data.ReportURL, len(data.Photos)),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	// Autonomous progression: Trigger Assessment Agent
	return uc.RunAssessment(ctx, claimID)
}
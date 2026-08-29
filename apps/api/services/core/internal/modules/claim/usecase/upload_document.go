package usecase

import (
	"context"
	"fmt"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) UploadDocument(ctx context.Context, claimID int, data *domain.RequestUploadDocument) (res domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:UploadDocument")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	if data.DocumentType == "" {
		data.DocumentType = "POLICE_REPORT"
	}
	if data.FileName == "" {
		data.FileName = "police_report_incident_km42.pdf"
	}
	if data.FileURL == "" {
		data.FileURL = fmt.Sprintf("https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0042/%s", data.FileName)
	}

	extractedData := `{"report_id": "PR-2026-9912", "officer": "Insp. H. Santoso", "incident_date": "2026-08-28 14:30", "location": "Highway KM 42", "severity": "MEDIUM", "liability": "SINGLE_VEHICLE_ACCIDENT"}`

	doc := shareddomain.ClaimDocument{
		ClaimID:       claimID,
		DocumentType:  data.DocumentType,
		FileName:      data.FileName,
		FileURL:       data.FileURL,
		Status:        "VERIFIED",
		ExtractedData: extractedData,
		UploadedAt:    time.Now(),
	}

	err = uc.repoSQL.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repoSQL.ClaimRepo().AddDocument(txCtx, &doc); err != nil {
			return err
		}

		event := shareddomain.ClaimEvent{
			ClaimID:   claimID,
			ActorName: "Claims Officer (Stand-in)",
			ActorType: shareddomain.ActorOfficer,
			Action:    "DOCUMENT_UPLOADED",
			Payload:   fmt.Sprintf(`{"document_type":"%s","file_name":"%s"}`, doc.DocumentType, doc.FileName),
		}
		return uc.repoSQL.ClaimRepo().AddEvent(txCtx, &event)
	})

	if err != nil {
		return res, err
	}

	// Autonomous next step: Intake Agent re-evaluates completeness and triggers assignment
	return uc.EvaluateIntake(ctx, claimID)
}

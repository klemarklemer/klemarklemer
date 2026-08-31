package domain

import (
	"time"

	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// ResponseClaimList model
type ResponseClaimList struct {
	Meta candishared.Meta `json:"meta"`
	Data []ResponseClaim  `json:"data"`
}

// ResponseClaim model
type ResponseClaim struct {
	ID                   int                                    `json:"id"`
	ClaimNumber          string                                 `json:"claim_number"`
	PolicyID             int                                    `json:"policy_id"`
	Policy               *shareddomain.Policy                   `json:"policy,omitempty"`
	Stage                string                                 `json:"stage"`
	DocumentCompleteness string                                 `json:"document_completeness"`
	SurveyRequired       bool                                   `json:"survey_required"`
	SurveyorID           *int                                   `json:"surveyor_id,omitempty"`
	Surveyor             *shareddomain.ClaimsOfficer            `json:"surveyor,omitempty"`
	SurveyStatus         string                                 `json:"survey_status,omitempty"`
	SurveySLADueAt       *string                                `json:"survey_sla_due_at,omitempty"`
	SurveyCompletedAt    *string                                `json:"survey_completed_at,omitempty"`
	SurveyReportURL      string                                 `json:"survey_report_url,omitempty"`
	SurveyPhotos         []string                               `json:"survey_photos,omitempty"`
	ClaimType            string                                 `json:"claim_type"`
	Severity             string                                 `json:"severity"`
	IncidentDescription  string                                 `json:"incident_description"`
	EstimatedLoss        float64                                `json:"estimated_loss"`
	ApprovedAmount       float64                                `json:"approved_amount"`
	FraudSignal          *string                                `json:"fraud_signal,omitempty"`
	CurrentOfficerID     *int                                   `json:"current_officer_id,omitempty"`
	CurrentOfficer       *shareddomain.ClaimsOfficer            `json:"current_officer,omitempty"`
	ClaimSLADueAt        *string                                `json:"claim_sla_due_at,omitempty"`
	StageSLADueAt        *string                                `json:"stage_sla_due_at,omitempty"`
	Status               string                                 `json:"status"`
	CreatedAt            string                                 `json:"created_at"`
	UpdatedAt            string                                 `json:"updated_at"`
	Documents            []shareddomain.ClaimDocument           `json:"documents,omitempty"`
	Events               []shareddomain.ClaimEvent              `json:"events,omitempty"`
	Assignment           *shareddomain.Assignment               `json:"assignment,omitempty"`
	Recommendation       *shareddomain.AssessmentRecommendation `json:"recommendation,omitempty"`
}

// Serialize from db model
func (r *ResponseClaim) Serialize(source *shareddomain.Claim) {
	r.ID = source.ID
	r.ClaimNumber = source.ClaimNumber
	r.PolicyID = source.PolicyID
	r.Policy = source.Policy
	r.Stage = source.Stage
	r.DocumentCompleteness = source.DocumentCompleteness
	r.SurveyRequired = source.SurveyRequired
	r.SurveyorID = source.SurveyorID
	r.Surveyor = source.Surveyor
	r.SurveyStatus = source.SurveyStatus
	r.SurveyReportURL = source.SurveyReportURL
	r.SurveyPhotos = source.SurveyPhotos
	if source.SurveySLADueAt != nil {
		due := source.SurveySLADueAt.Format(time.RFC3339)
		r.SurveySLADueAt = &due
	}
	if source.SurveyCompletedAt != nil {
		done := source.SurveyCompletedAt.Format(time.RFC3339)
		r.SurveyCompletedAt = &done
	}
	r.ClaimType = source.ClaimType
	r.Severity = source.Severity
	r.IncidentDescription = source.IncidentDescription
	r.EstimatedLoss = source.EstimatedLoss
	r.ApprovedAmount = source.ApprovedAmount
	r.FraudSignal = source.FraudSignal
	r.CurrentOfficerID = source.CurrentOfficerID
	r.CurrentOfficer = source.CurrentOfficer
	r.Status = source.Status
	r.CreatedAt = source.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = source.UpdatedAt.Format(time.RFC3339)

	if source.ClaimSLADueAt != nil {
		sla := source.ClaimSLADueAt.Format(time.RFC3339)
		r.ClaimSLADueAt = &sla
	}
	if source.StageSLADueAt != nil {
		sla := source.StageSLADueAt.Format(time.RFC3339)
		r.StageSLADueAt = &sla
	}

	r.Documents = source.Documents
	r.Events = source.Events
	r.Assignment = source.Assignment
	r.Recommendation = source.Recommendation
}

package domain

import (
	shareddomain "monorepo/services/core/pkg/shared/domain"
)

// RequestCreateClaim model
type RequestCreateClaim struct {
	PolicyID            int     `json:"policy_id"`
	ClaimType           string  `json:"claim_type"`
	IncidentDescription string  `json:"incident_description"`
	EstimatedLoss       float64 `json:"estimated_loss"`
}

// RequestUploadDocument model
type RequestUploadDocument struct {
	DocumentType string `json:"document_type"`
	FileName     string `json:"file_name"`
	FileURL      string `json:"file_url"`
	RawContent   string `json:"raw_content"`
}

// RequestHumanApproval model
type RequestHumanApproval struct {
	OfficerID int    `json:"officer_id"`
	Action    string `json:"action"` // APPROVE, REJECT
	Notes     string `json:"notes"`
}

// RequestCompleteSurvey model
type RequestCompleteSurvey struct {
	SurveyorID     int      `json:"surveyor_id"`
	ReportURL      string   `json:"report_url"`
	Photos         []string `json:"photos"`
	Notes          string   `json:"notes"`
}

// RequestClaim placeholder for candi compatibility
type RequestClaim struct {
	ID                  int     `json:"id"`
	PolicyID            int     `json:"policy_id"`
	ClaimType           string  `json:"claim_type"`
	IncidentDescription string  `json:"incident_description"`
	EstimatedLoss       float64 `json:"estimated_loss"`
}

func (r *RequestClaim) Deserialize() (res shareddomain.Claim) {
	res.ID = r.ID
	res.PolicyID = r.PolicyID
	res.ClaimType = r.ClaimType
	res.IncidentDescription = r.IncidentDescription
	res.EstimatedLoss = r.EstimatedLoss
	return
}

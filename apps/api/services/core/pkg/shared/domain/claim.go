package domain

import (
	"time"
)

// Claim stage constants strictly adhering to CONTEXT.md glossary
const (
	StageIntake               = "INTAKE"
	StageDocumentVerification = "DOCUMENT_VERIFICATION"
	StageAssignment           = "ASSIGNMENT"
	StageAssessment           = "ASSESSMENT"
	StageDecision             = "DECISION"
	StageClosed               = "CLOSED"
)

// Document completeness constants
const (
	CompletenessComplete   = "COMPLETE"
	CompletenessIncomplete = "INCOMPLETE"
)

// Recommendation outcomes
const (
	OutcomeApprove      = "APPROVE"
	OutcomeReject       = "REJECT"
	OutcomeManualReview = "MANUAL_REVIEW"
)

// Actor types for immutable Claim events
const (
	ActorAgent   = "AGENT"
	ActorOfficer = "OFFICER"
	ActorSystem  = "SYSTEM"
)

// Claim model
type Claim struct {
	ID                   int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClaimNumber          string    `gorm:"column:claim_number;type:varchar(64);uniqueIndex;not null" json:"claim_number"`
	PolicyID             int       `gorm:"column:policy_id;not null" json:"policy_id"`
	Policy               *Policy   `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Stage                string    `gorm:"column:stage;type:varchar(64);not null;default:'INTAKE'" json:"stage"`
	DocumentCompleteness string    `gorm:"column:document_completeness;type:varchar(32);not null;default:'INCOMPLETE'" json:"document_completeness"`
	SurveyRequired       bool      `gorm:"column:survey_required;not null;default:false" json:"survey_required"`
	ClaimType            string    `gorm:"column:claim_type;type:varchar(64);not null;default:'MOTOR'" json:"claim_type"`
	Severity             string    `gorm:"column:severity;type:varchar(32);not null;default:'MEDIUM'" json:"severity"`
	IncidentDescription  string    `gorm:"column:incident_description;type:text" json:"incident_description"`
	EstimatedLoss        float64   `gorm:"column:estimated_loss;type:decimal(15,2);not null;default:0" json:"estimated_loss"`
	ApprovedAmount       float64   `gorm:"column:approved_amount;type:decimal(15,2);not null;default:0" json:"approved_amount"`
	FraudSignal          *string    `gorm:"column:fraud_signal;type:varchar(512)" json:"fraud_signal,omitempty"`
	CurrentOfficerID     *int      `gorm:"column:current_officer_id" json:"current_officer_id,omitempty"`
	CurrentOfficer       *ClaimsOfficer `gorm:"foreignKey:CurrentOfficerID" json:"current_officer,omitempty"`
	ClaimSLADueAt        *time.Time `gorm:"column:claim_sla_due_at" json:"claim_sla_due_at,omitempty"`
	StageSLADueAt        *time.Time `gorm:"column:stage_sla_due_at" json:"stage_sla_due_at,omitempty"`
	Status               string    `gorm:"column:status;type:varchar(32);not null;default:'OPEN'" json:"status"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Documents            []ClaimDocument           `gorm:"foreignKey:ClaimID" json:"documents,omitempty"`
	Events               []ClaimEvent              `gorm:"foreignKey:ClaimID" json:"events,omitempty"`
	Assignment           *Assignment               `gorm:"foreignKey:ClaimID" json:"assignment,omitempty"`
	Recommendation       *AssessmentRecommendation `gorm:"foreignKey:ClaimID" json:"recommendation,omitempty"`
}

// TableName return table name of Claim model
func (Claim) TableName() string {
	return "claims"
}

// ClaimDocument model representing required artifacts for completeness
type ClaimDocument struct {
	ID            int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClaimID       int       `gorm:"column:claim_id;not null;index" json:"claim_id"`
	DocumentType  string    `gorm:"column:document_type;type:varchar(64);not null" json:"document_type"`
	FileName      string    `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FileURL       string    `gorm:"column:file_url;type:text;not null" json:"file_url"`
	Status        string    `gorm:"column:status;type:varchar(32);not null;default:'VERIFIED'" json:"status"`
	ExtractedData string    `gorm:"column:extracted_data;type:text" json:"extracted_data"`
	UploadedAt    time.Time `gorm:"column:uploaded_at;autoCreateTime" json:"uploaded_at"`
}

func (ClaimDocument) TableName() string {
	return "claim_documents"
}

// ClaimEvent model for immutable timeline facts
type ClaimEvent struct {
	ID            int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClaimID       int       `gorm:"column:claim_id;not null;index" json:"claim_id"`
	ActorName     string    `gorm:"column:actor_name;type:varchar(128);not null" json:"actor_name"`
	ActorType     string    `gorm:"column:actor_type;type:varchar(32);not null" json:"actor_type"`
	Action        string    `gorm:"column:action;type:varchar(128);not null" json:"action"`
	PreviousStage string    `gorm:"column:previous_stage;type:varchar(64)" json:"previous_stage"`
	NewStage      string    `gorm:"column:new_stage;type:varchar(64)" json:"new_stage"`
	Payload       string    `gorm:"column:payload;type:text" json:"payload"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (ClaimEvent) TableName() string {
	return "claim_events"
}

// Assignment model for binding a claims officer to a claim
type Assignment struct {
	ID            int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClaimID       int       `gorm:"column:claim_id;uniqueIndex;not null" json:"claim_id"`
	OfficerID     int       `gorm:"column:officer_id;not null" json:"officer_id"`
	Officer       *ClaimsOfficer `gorm:"foreignKey:OfficerID" json:"officer,omitempty"`
	WorkloadScore float64   `gorm:"column:workload_score;type:decimal(5,2);not null" json:"workload_score"`
	SkillScore    float64   `gorm:"column:skill_score;type:decimal(5,2);not null" json:"skill_score"`
	TotalScore    float64   `gorm:"column:total_score;type:decimal(5,2);not null" json:"total_score"`
	AssignedAt    time.Time `gorm:"column:assigned_at;autoCreateTime" json:"assigned_at"`
}

func (Assignment) TableName() string {
	return "assignments"
}

// AssessmentRecommendation model for AI proposed outcome (never auto-approves)
type AssessmentRecommendation struct {
	ID          int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClaimID     int       `gorm:"column:claim_id;uniqueIndex;not null" json:"claim_id"`
	Outcome     string    `gorm:"column:outcome;type:varchar(32);not null" json:"outcome"`
	Confidence  float64   `gorm:"column:confidence;type:decimal(5,2);not null" json:"confidence"`
	Reasons     string    `gorm:"column:reasons;type:text;not null" json:"reasons"`
	GeneratedAt time.Time `gorm:"column:generated_at;autoCreateTime" json:"generated_at"`
}

func (AssessmentRecommendation) TableName() string {
	return "assessment_recommendations"
}

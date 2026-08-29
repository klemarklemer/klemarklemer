package repository

import (
	"context"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// ClaimRepository abstract interface
type ClaimRepository interface {
	FetchAll(ctx context.Context, filter *domain.FilterClaim) ([]shareddomain.Claim, error)
	Count(ctx context.Context, filter *domain.FilterClaim) int
	Find(ctx context.Context, filter *domain.FilterClaim) (shareddomain.Claim, error)
	Save(ctx context.Context, data *shareddomain.Claim, updateOptions ...candishared.DBUpdateOptionFunc) error
	Delete(ctx context.Context, filter *domain.FilterClaim) (err error)

	// Event operations
	AddEvent(ctx context.Context, event *shareddomain.ClaimEvent) error
	GetEventsByClaimID(ctx context.Context, claimID int) ([]shareddomain.ClaimEvent, error)

	// Document operations
	AddDocument(ctx context.Context, doc *shareddomain.ClaimDocument) error
	GetDocumentsByClaimID(ctx context.Context, claimID int) ([]shareddomain.ClaimDocument, error)

	// Assignment operations
	SaveAssignment(ctx context.Context, assign *shareddomain.Assignment) error
	GetAssignmentByClaimID(ctx context.Context, claimID int) (*shareddomain.Assignment, error)

	// Assessment recommendation operations
	SaveRecommendation(ctx context.Context, rec *shareddomain.AssessmentRecommendation) error
	GetRecommendationByClaimID(ctx context.Context, claimID int) (*shareddomain.AssessmentRecommendation, error)
}

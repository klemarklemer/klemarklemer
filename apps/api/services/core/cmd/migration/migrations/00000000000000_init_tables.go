package migrations

import (
	"monorepo/services/core/pkg/shared/domain"
)

// GetMigrateTables get migrate table list
func GetMigrateTables() []any {
	return []any{
		&domain.Policy{},
		&domain.ClaimsOfficer{},
		&domain.Claim{},
		&domain.ClaimDocument{},
		&domain.ClaimEvent{},
		&domain.Assignment{},
		&domain.AssessmentRecommendation{},
	}
}

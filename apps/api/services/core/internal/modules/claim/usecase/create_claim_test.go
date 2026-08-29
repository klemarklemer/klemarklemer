package usecase

import (
	"context"
	"testing"

	"monorepo/services/core/internal/modules/claim/domain"

	"github.com/stretchr/testify/assert"
)

func Test_claimUsecaseImpl_CreateClaim(t *testing.T) {
	t.Run("Testcase #1: Positive - requires DB, skip in unit test", func(t *testing.T) {
		// Integration test: CreateClaim requires a database connection.
		// Skipped in unit test; covered by integration test with running Postgres.
		t.Skip("requires database connection")
		uc := claimUsecaseImpl{}
		_, err := uc.CreateClaim(context.Background(), &domain.RequestCreateClaim{
			PolicyID:            1,
			ClaimType:           "MOTOR",
			IncidentDescription: "Test incident",
			EstimatedLoss:       1000.00,
		})
		assert.NoError(t, err)
	})
}

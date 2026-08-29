package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/policy/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *policyUsecaseImpl) UpdatePolicy(ctx context.Context, data *domain.RequestPolicy) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyUsecase:UpdatePolicy")
	defer trace.Finish()

	repoFilter := domain.FilterPolicy{ID: &data.ID}
	existing, err := uc.repoSQL.PolicyRepo().Find(ctx, &repoFilter)
	if err != nil {
		return err
	}
	existing.PolicyHolderName = data.PolicyHolderName
	existing.Status = data.Status
	existing.MaxCoverageAmount = data.MaxCoverageAmount
	existing.DeductibleAmount = data.DeductibleAmount
	err = uc.repoSQL.WithTransaction(ctx, func(ctx context.Context) error {
		return uc.repoSQL.PolicyRepo().Save(ctx, &existing, candishared.DBUpdateSetUpdatedFields("PolicyHolderName", "Status", "MaxCoverageAmount", "DeductibleAmount"))
	})
	return
}

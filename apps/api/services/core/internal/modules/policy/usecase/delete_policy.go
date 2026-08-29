package usecase

import (
	"context"
	
	"monorepo/services/core/internal/modules/policy/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *policyUsecaseImpl) DeletePolicy(ctx context.Context, id int) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyUsecase:DeletePolicy")
	defer trace.Finish()

	repoFilter := domain.FilterPolicy{ID: &id}
	return uc.repoSQL.PolicyRepo().Delete(ctx, &repoFilter)
}

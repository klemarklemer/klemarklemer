package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/policy/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *policyUsecaseImpl) GetDetailPolicy(ctx context.Context, id int) (result domain.ResponsePolicy, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyUsecase:GetDetailPolicy")
	defer trace.Finish()

	repoFilter := domain.FilterPolicy{ID: &id}
	data, err := uc.repoSQL.PolicyRepo().Find(ctx, &repoFilter)
	if err != nil {
		return result, err
	}

	result.Serialize(&data)
	return
}

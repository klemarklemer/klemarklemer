package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/policy/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *policyUsecaseImpl) GetAllPolicy(ctx context.Context, filter *domain.FilterPolicy) (result domain.ResponsePolicyList, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyUsecase:GetAllPolicy")
	defer trace.Finish()

	data, err := uc.repoSQL.PolicyRepo().FetchAll(ctx, filter)
	if err != nil {
		return result, err
	}
	count := uc.repoSQL.PolicyRepo().Count(ctx, filter)
	result.Meta = candishared.NewMeta(filter.Page, filter.Limit, count)

	result.Data = make([]domain.ResponsePolicy, len(data))
	for i, detail := range data {
		result.Data[i].Serialize(&detail)
	}

	return
}

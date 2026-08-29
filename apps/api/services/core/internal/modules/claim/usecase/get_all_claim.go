package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/claim/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) GetAllClaim(ctx context.Context, filter *domain.FilterClaim) (data domain.ResponseClaimList, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:GetAllClaim")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	count := uc.repoSQL.ClaimRepo().Count(ctx, filter)
	data.Meta = candishared.NewMeta(filter.Page, filter.Limit, count)

	claims, err := uc.repoSQL.ClaimRepo().FetchAll(ctx, filter)
	if err != nil {
		return data, err
	}

	data.Data = make([]domain.ResponseClaim, len(claims))
	for i := range claims {
		data.Data[i].Serialize(&claims[i])
	}
	return data, nil
}

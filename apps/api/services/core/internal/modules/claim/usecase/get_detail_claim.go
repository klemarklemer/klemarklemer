package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/claim/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *claimUsecaseImpl) GetDetailClaim(ctx context.Context, id int) (data domain.ResponseClaim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimUsecase:GetDetailClaim")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	filter := domain.FilterClaim{ID: &id}
	claim, err := uc.repoSQL.ClaimRepo().Find(ctx, &filter)
	if err != nil {
		return data, err
	}

	data.Serialize(&claim)
	return data, nil
}

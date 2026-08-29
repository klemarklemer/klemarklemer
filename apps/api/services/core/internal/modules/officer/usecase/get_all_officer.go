package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/officer/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *officerUsecaseImpl) GetAllOfficer(ctx context.Context, filter *domain.FilterOfficer) (result domain.ResponseOfficerList, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerUsecase:GetAllOfficer")
	defer trace.Finish()

	data, err := uc.repoSQL.OfficerRepo().FetchAll(ctx, filter)
	if err != nil {
		return result, err
	}
	count := uc.repoSQL.OfficerRepo().Count(ctx, filter)
	result.Meta = candishared.NewMeta(filter.Page, filter.Limit, count)

	result.Data = make([]domain.ResponseOfficer, len(data))
	for i, detail := range data {
		result.Data[i].Serialize(&detail)
	}

	return
}

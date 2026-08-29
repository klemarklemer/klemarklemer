package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/officer/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *officerUsecaseImpl) GetDetailOfficer(ctx context.Context, id int) (result domain.ResponseOfficer, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerUsecase:GetDetailOfficer")
	defer trace.Finish()

	repoFilter := domain.FilterOfficer{ID: &id}
	data, err := uc.repoSQL.OfficerRepo().Find(ctx, &repoFilter)
	if err != nil {
		return result, err
	}

	result.Serialize(&data)
	return
}

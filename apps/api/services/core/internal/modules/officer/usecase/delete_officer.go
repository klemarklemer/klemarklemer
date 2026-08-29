package usecase

import (
	"context"
	
	"monorepo/services/core/internal/modules/officer/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *officerUsecaseImpl) DeleteOfficer(ctx context.Context, id int) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerUsecase:DeleteOfficer")
	defer trace.Finish()

	repoFilter := domain.FilterOfficer{ID: &id}
	return uc.repoSQL.OfficerRepo().Delete(ctx, &repoFilter)
}

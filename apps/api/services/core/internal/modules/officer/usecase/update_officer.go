package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/officer/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *officerUsecaseImpl) UpdateOfficer(ctx context.Context, data *domain.RequestOfficer) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerUsecase:UpdateOfficer")
	defer trace.Finish()

	repoFilter := domain.FilterOfficer{ID: &data.ID}
	existing, err := uc.repoSQL.OfficerRepo().Find(ctx, &repoFilter)
	if err != nil {
		return err
	}
	existing.Name = data.Name
	existing.Email = data.Email
	existing.Role = data.Role
	existing.IsAvailable = data.IsAvailable
	existing.MotorSkillRating = data.MotorSkillRating
	err = uc.repoSQL.WithTransaction(ctx, func(ctx context.Context) error {
		return uc.repoSQL.OfficerRepo().Save(ctx, &existing, candishared.DBUpdateSetUpdatedFields("Name", "Email", "Role", "IsAvailable", "MotorSkillRating"))
	})
	return
}

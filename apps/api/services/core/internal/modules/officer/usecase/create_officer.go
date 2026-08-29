package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/officer/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *officerUsecaseImpl) CreateOfficer(ctx context.Context, req *domain.RequestOfficer) (result domain.ResponseOfficer, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerUsecase:CreateOfficer")
	defer trace.Finish()

	data := req.Deserialize()
	err = uc.repoSQL.OfficerRepo().Save(ctx, &data)
	result.Serialize(&data)

	// Sample using broker publisher
	// uc.deps.GetBroker(types.Kafka). // get registered broker type (sample Kafka)
	// 				GetPublisher().
	// 				PublishMessage(ctx, &candishared.PublisherArgument{
	// 		Topic:   "[topic]",
	// 		Key:     "[key]",
	// 		Message: candihelper.ToBytes("[message]"),
	// 	})
	return
}

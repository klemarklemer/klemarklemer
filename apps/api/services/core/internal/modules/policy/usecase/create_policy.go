package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/policy/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *policyUsecaseImpl) CreatePolicy(ctx context.Context, req *domain.RequestPolicy) (result domain.ResponsePolicy, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyUsecase:CreatePolicy")
	defer trace.Finish()

	data := req.Deserialize()
	err = uc.repoSQL.PolicyRepo().Save(ctx, &data)
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

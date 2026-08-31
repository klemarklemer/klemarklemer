package usecase

import (
	"context"

	"monorepo/services/core/internal/modules/claim/domain"
	"monorepo/services/core/pkg/shared/gemini"
	"monorepo/services/core/pkg/shared/repository"
	"monorepo/services/core/pkg/shared/usecase/common"

	"github.com/golangid/candi/codebase/factory/dependency"
	"github.com/golangid/candi/logger"
)

// ClaimUsecase abstraction
type ClaimUsecase interface {
	GetAllClaim(ctx context.Context, filter *domain.FilterClaim) (data domain.ResponseClaimList, err error)
	GetDetailClaim(ctx context.Context, id int) (data domain.ResponseClaim, err error)
	CreateClaim(ctx context.Context, data *domain.RequestCreateClaim) (res domain.ResponseClaim, err error)
	UploadDocument(ctx context.Context, claimID int, data *domain.RequestUploadDocument) (res domain.ResponseClaim, err error)
	EvaluateIntake(ctx context.Context, claimID int) (res domain.ResponseClaim, err error)
	RunAssignment(ctx context.Context, claimID int) (res domain.ResponseClaim, err error)
	RunAssessment(ctx context.Context, claimID int) (res domain.ResponseClaim, err error)
	SubmitHumanApproval(ctx context.Context, claimID int, data *domain.RequestHumanApproval) (res domain.ResponseClaim, err error)
	ResetDemo(ctx context.Context) (res domain.ResponseClaim, err error)

	// Survey operations
	AssignSurveyor(ctx context.Context, claimID int, surveyorID int) (res domain.ResponseClaim, err error)
	UpdateSurveyStatus(ctx context.Context, claimID int, status string) (res domain.ResponseClaim, err error)
	CompleteSurvey(ctx context.Context, claimID int, data *domain.RequestCompleteSurvey) (res domain.ResponseClaim, err error)

	UpdateClaim(ctx context.Context, data *domain.RequestClaim) (err error)
	DeleteClaim(ctx context.Context, id int) (err error)
}

type claimUsecaseImpl struct {
	deps          dependency.Dependency
	sharedUsecase common.Usecase
	repoSQL       repository.RepoSQL
	// classifier backs Loop 1 and assessor backs Loop 3. Both are interfaces so
	// the reasoning can be tested without a Gemini credential, and so an
	// unconfigured deployment falls back to deterministic rules.
	classifier gemini.Classifier
	assessor   gemini.Assessor
}

// NewClaimUsecase constructor
func NewClaimUsecase(deps dependency.Dependency) (ClaimUsecase, func(sharedUsecase common.Usecase)) {
	assessor, assessEngine := gemini.New(context.Background())
	classifier, intakeEngine := gemini.NewClassifier(context.Background())
	logger.LogGreen("claim usecase: intake agent engine -> " + intakeEngine)
	logger.LogGreen("claim usecase: assessment agent engine -> " + assessEngine)

	uc := &claimUsecaseImpl{
		deps:       deps,
		repoSQL:    repository.GetSharedRepoSQL(),
		classifier: classifier,
		assessor:   assessor,
	}
	return uc, func(sharedUsecase common.Usecase) {
		uc.sharedUsecase = sharedUsecase
	}
}

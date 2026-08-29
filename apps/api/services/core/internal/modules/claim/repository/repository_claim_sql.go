package repository

import (
	"context"
	"strings"
	"time"

	"monorepo/services/core/internal/modules/claim/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"monorepo/globalshared"
)

type claimRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewClaimRepoSQL GORM SQL repository constructor
func NewClaimRepoSQL(readDB, writeDB *gorm.DB) ClaimRepository {
	return &claimRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *claimRepoSQL) getWriteDB(ctx context.Context) *gorm.DB {
	if tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB); ok {
		return tx
	}
	return r.writeDB
}

func (r *claimRepoSQL) FetchAll(ctx context.Context, filter *domain.FilterClaim) (data []shareddomain.Claim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:FetchAll")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	if filter.OrderBy == "" {
		filter.OrderBy = "id"
	}

	db := r.setFilterClaim(globalshared.SetSpanToGorm(ctx, r.readDB), filter).Order(clause.OrderByColumn{
		Column: clause.Column{Name: filter.OrderBy},
		Desc:   strings.ToUpper(filter.Sort) == "DESC",
	})
	if filter.Limit > 0 || !filter.ShowAll {
		db = db.Limit(filter.Limit).Offset(filter.CalculateOffset())
	}
	err = db.Preload("Policy").Preload("CurrentOfficer").Preload("Documents").Preload("Events", func(db *gorm.DB) *gorm.DB {
		return db.Order("claim_events.created_at ASC")
	}).Preload("Assignment").Preload("Assignment.Officer").Preload("Recommendation").Find(&data).Error
	return
}

func (r *claimRepoSQL) Count(ctx context.Context, filter *domain.FilterClaim) (count int) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:Count")
	defer trace.Finish()

	var total int64
	r.setFilterClaim(globalshared.SetSpanToGorm(ctx, r.readDB), filter).Model(&shareddomain.Claim{}).Count(&total)
	count = int(total)
	trace.Log("count", count)
	return
}

func (r *claimRepoSQL) Find(ctx context.Context, filter *domain.FilterClaim) (result shareddomain.Claim, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterClaim(globalshared.SetSpanToGorm(ctx, r.readDB), filter).
		Preload("Policy").
		Preload("CurrentOfficer").
		Preload("Documents").
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Order("claim_events.created_at ASC")
		}).
		Preload("Assignment").
		Preload("Assignment.Officer").
		Preload("Recommendation").
		First(&result).Error
	return
}

func (r *claimRepoSQL) Save(ctx context.Context, data *shareddomain.Claim, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:Save")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	data.UpdatedAt = time.Now()
	if data.CreatedAt.IsZero() {
		data.CreatedAt = time.Now()
	}
	if data.ID == 0 {
		err = globalshared.SetSpanToGorm(ctx, db).Omit(clause.Associations).Create(data).Error
	} else {
		err = globalshared.SetSpanToGorm(ctx, db).Model(data).Omit(clause.Associations).Updates(r.updateTools.ToMap(data, updateOptions...)).Error
	}
	return
}

func (r *claimRepoSQL) Delete(ctx context.Context, filter *domain.FilterClaim) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:Delete")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	err = r.setFilterClaim(globalshared.SetSpanToGorm(ctx, db), filter).Delete(&shareddomain.Claim{}).Error
	return
}

func (r *claimRepoSQL) AddEvent(ctx context.Context, event *shareddomain.ClaimEvent) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:AddEvent")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	err = globalshared.SetSpanToGorm(ctx, db).Create(event).Error
	return
}

func (r *claimRepoSQL) GetEventsByClaimID(ctx context.Context, claimID int) (events []shareddomain.ClaimEvent, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:GetEventsByClaimID")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = globalshared.SetSpanToGorm(ctx, r.readDB).
		Where("claim_id = ?", claimID).
		Order("created_at ASC").
		Find(&events).Error
	return
}

func (r *claimRepoSQL) AddDocument(ctx context.Context, doc *shareddomain.ClaimDocument) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:AddDocument")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	if doc.UploadedAt.IsZero() {
		doc.UploadedAt = time.Now()
	}
	err = globalshared.SetSpanToGorm(ctx, db).Create(doc).Error
	return
}

func (r *claimRepoSQL) GetDocumentsByClaimID(ctx context.Context, claimID int) (docs []shareddomain.ClaimDocument, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:GetDocumentsByClaimID")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = globalshared.SetSpanToGorm(ctx, r.readDB).
		Where("claim_id = ?", claimID).
		Order("uploaded_at ASC").
		Find(&docs).Error
	return
}

func (r *claimRepoSQL) SaveAssignment(ctx context.Context, assign *shareddomain.Assignment) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:SaveAssignment")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	if assign.AssignedAt.IsZero() {
		assign.AssignedAt = time.Now()
	}
	err = globalshared.SetSpanToGorm(ctx, db).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "claim_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"officer_id", "workload_score", "skill_score", "total_score", "assigned_at"}),
		}).
		Create(assign).Error
	return
}

func (r *claimRepoSQL) GetAssignmentByClaimID(ctx context.Context, claimID int) (assign *shareddomain.Assignment, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:GetAssignmentByClaimID")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	var data shareddomain.Assignment
	err = globalshared.SetSpanToGorm(ctx, r.readDB).
		Preload("Officer").
		Where("claim_id = ?", claimID).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *claimRepoSQL) SaveRecommendation(ctx context.Context, rec *shareddomain.AssessmentRecommendation) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:SaveRecommendation")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	if rec.GeneratedAt.IsZero() {
		rec.GeneratedAt = time.Now()
	}
	err = globalshared.SetSpanToGorm(ctx, db).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "claim_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"outcome", "confidence", "reasons", "generated_at"}),
		}).
		Create(rec).Error
	return
}

func (r *claimRepoSQL) GetRecommendationByClaimID(ctx context.Context, claimID int) (rec *shareddomain.AssessmentRecommendation, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "ClaimRepoSQL:GetRecommendationByClaimID")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	var data shareddomain.AssessmentRecommendation
	err = globalshared.SetSpanToGorm(ctx, r.readDB).
		Where("claim_id = ?", claimID).
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *claimRepoSQL) setFilterClaim(db *gorm.DB, filter *domain.FilterClaim) *gorm.DB {
	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}
	if filter.ClaimNumber != "" {
		db = db.Where("claim_number = ?", filter.ClaimNumber)
	}
	if filter.Stage != "" {
		db = db.Where("stage = ?", filter.Stage)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		db = db.Where("(claim_number ILIKE '%%' || ? || '%%' OR incident_description ILIKE '%%' || ? || '%%')", filter.Search, filter.Search)
	}

	for _, preload := range filter.Preloads {
		db = db.Preload(preload)
	}

	return db
}

package repository

import (
	"context"
	"strings"
	"time"

	"monorepo/services/core/internal/modules/policy/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"monorepo/globalshared"
)

type policyRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewPolicyRepoSQL GORM SQL repository constructor
func NewPolicyRepoSQL(readDB, writeDB *gorm.DB) PolicyRepository {
	return &policyRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *policyRepoSQL) getWriteDB(ctx context.Context) *gorm.DB {
	if tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB); ok {
		return tx
	}
	return r.writeDB
}

func (r *policyRepoSQL) FetchAll(ctx context.Context, filter *domain.FilterPolicy) (data []shareddomain.Policy, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyRepoSQL:FetchAll")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	if filter.OrderBy == "" {
		filter.OrderBy = "id"
	}

	db := r.setFilterPolicy(globalshared.SetSpanToGorm(ctx, r.readDB), filter).Order(clause.OrderByColumn{
		Column: clause.Column{Name: filter.OrderBy},
		Desc:   strings.ToUpper(filter.Sort) == "DESC",
	})
	if filter.Limit > 0 || !filter.ShowAll {
		db = db.Limit(filter.Limit).Offset(filter.CalculateOffset())
	}
	err = db.Find(&data).Error
	return
}

func (r *policyRepoSQL) Count(ctx context.Context, filter *domain.FilterPolicy) (count int) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyRepoSQL:Count")
	defer trace.Finish()

	var total int64
	r.setFilterPolicy(globalshared.SetSpanToGorm(ctx, r.readDB), filter).Model(&shareddomain.Policy{}).Count(&total)
	count = int(total)
	trace.Log("count", count)
	return
}

func (r *policyRepoSQL) Find(ctx context.Context, filter *domain.FilterPolicy) (result shareddomain.Policy, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterPolicy(globalshared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *policyRepoSQL) Save(ctx context.Context, data *shareddomain.Policy, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyRepoSQL:Save")
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

func (r *policyRepoSQL) Delete(ctx context.Context, filter *domain.FilterPolicy) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PolicyRepoSQL:Delete")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	err = r.setFilterPolicy(globalshared.SetSpanToGorm(ctx, db), filter).Delete(&shareddomain.Policy{}).Error
	return
}

func (r *policyRepoSQL) setFilterPolicy(db *gorm.DB, filter *domain.FilterPolicy) *gorm.DB {
	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}
	if filter.Search != "" {
		db = db.Where("(policy_number ILIKE '%%' || ? || '%%' OR policy_holder_name ILIKE '%%' || ? || '%%' OR vehicle_plate ILIKE '%%' || ? || '%%')", filter.Search, filter.Search, filter.Search)
	}

	for _, preload := range filter.Preloads {
		db = db.Preload(preload)
	}

	return db
}

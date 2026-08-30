package repository

import (
	"context"
	"strings"
	"time"

	"monorepo/services/core/internal/modules/officer/domain"
	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"monorepo/globalshared"
)

type officerRepoSQL struct {
	readDB, writeDB *gorm.DB
	updateTools     *candishared.DBUpdateTools
}

// NewOfficerRepoSQL GORM SQL repository constructor
func NewOfficerRepoSQL(readDB, writeDB *gorm.DB) OfficerRepository {
	return &officerRepoSQL{
		readDB: readDB, writeDB: writeDB,
		updateTools: &candishared.DBUpdateTools{
			KeyExtractorFunc: candishared.DBUpdateGORMExtractorKey, IgnoredFields: []string{"id"},
		},
	}
}

func (r *officerRepoSQL) getWriteDB(ctx context.Context) *gorm.DB {
	if tx, ok := candishared.GetValueFromContext(ctx, candishared.ContextKeySQLTransaction).(*gorm.DB); ok {
		return tx
	}
	return r.writeDB
}

func (r *officerRepoSQL) FetchAll(ctx context.Context, filter *domain.FilterOfficer) (data []shareddomain.Officer, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerRepoSQL:FetchAll")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	if filter.OrderBy == "" {
		filter.OrderBy = "id"
	}

	db := r.setFilterOfficer(globalshared.SetSpanToGorm(ctx, r.readDB), filter).Order(clause.OrderByColumn{
		Column: clause.Column{Name: filter.OrderBy},
		Desc:   strings.ToUpper(filter.Sort) == "DESC",
	})
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit).Offset(filter.CalculateOffset())
	}
	err = db.Find(&data).Error
	return
}

func (r *officerRepoSQL) Count(ctx context.Context, filter *domain.FilterOfficer) (count int) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerRepoSQL:Count")
	defer trace.Finish()

	var total int64
	r.setFilterOfficer(globalshared.SetSpanToGorm(ctx, r.readDB), filter).Model(&shareddomain.Officer{}).Count(&total)
	count = int(total)
	trace.Log("count", count)
	return
}

func (r *officerRepoSQL) Find(ctx context.Context, filter *domain.FilterOfficer) (result shareddomain.Officer, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerRepoSQL:Find")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	err = r.setFilterOfficer(globalshared.SetSpanToGorm(ctx, r.readDB), filter).First(&result).Error
	return
}

func (r *officerRepoSQL) Save(ctx context.Context, data *shareddomain.Officer, updateOptions ...candishared.DBUpdateOptionFunc) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerRepoSQL:Save")
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

func (r *officerRepoSQL) Delete(ctx context.Context, filter *domain.FilterOfficer) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "OfficerRepoSQL:Delete")
	defer func() { trace.Finish(tracer.FinishWithError(err)) }()

	db := r.getWriteDB(ctx)
	err = r.setFilterOfficer(globalshared.SetSpanToGorm(ctx, db), filter).Delete(&shareddomain.Officer{}).Error
	return
}

func (r *officerRepoSQL) setFilterOfficer(db *gorm.DB, filter *domain.FilterOfficer) *gorm.DB {
	if filter.ID != nil {
		db = db.Where("id = ?", *filter.ID)
	}
	if filter.Search != "" {
		db = db.Where("(name ILIKE '%%' || ? || '%%' OR email ILIKE '%%' || ? || '%%')", filter.Search, filter.Search)
	}

	for _, preload := range filter.Preloads {
		db = db.Preload(preload)
	}

	return db
}

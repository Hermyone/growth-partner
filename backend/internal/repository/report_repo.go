// growth-partner/backend/internal/repository/report_repo.go
// 报告仓库接口和实现

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type ReportRepository interface {
	Create(ctx context.Context, report *model.Report) error
	FindByID(ctx context.Context, id uint64) (*model.Report, error)
	FindByClassID(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Report, int64, error)
	Update(ctx context.Context, report *model.Report) error
	UpdateStatus(ctx context.Context, id uint64, status model.ReportStatus, fileURL string, fileSize int64) error
	Delete(ctx context.Context, id uint64) error
}

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) Create(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.Report, error) {
	var report model.Report
	err := r.db.WithContext(ctx).First(&report, id).Error
	return &report, err
}

func (r *reportRepositoryImpl) FindByClassID(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Report, int64, error) {
	var reports []*model.Report
	db := r.db.WithContext(ctx).Where("class_id = ?", classID)

	// 应用过滤条件
	if reportType, ok := params["report_type"].(string); ok && reportType != "" {
		db = db.Where("report_type = ?", reportType)
	}
	if status, ok := params["status"].(string); ok && status != "" {
		db = db.Where("status = ?", status)
	}
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		db = db.Where("start_date >= ?", startDate)
	}
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		db = db.Where("end_date <= ?", endDate)
	}

	// 分页
	var count int64
	if err := db.Model(&model.Report{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "created_at DESC"
	if o, ok := params["order"].(string); ok && o != "" {
		order = o
	}
	db = db.Order(order)

	// 分页
	if page, ok := params["page"].(int); ok && page > 0 {
		limit := 10
		if l, ok := params["limit"].(int); ok && l > 0 {
			limit = l
		}
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	err := db.Find(&reports).Error
	return reports, count, err
}

func (r *reportRepositoryImpl) Update(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepositoryImpl) UpdateStatus(ctx context.Context, id uint64, status model.ReportStatus, fileURL string, fileSize int64) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if fileURL != "" {
		updates["file_url"] = fileURL
		updates["file_size"] = fileSize
	}
	return r.db.WithContext(ctx).Model(&model.Report{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reportRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Report{}, id).Error
}

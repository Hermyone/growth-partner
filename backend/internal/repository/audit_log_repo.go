// growth-partner/backend/internal/repository/audit_log_repo.go
// 审计日志仓库接口和实现

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	FindByID(ctx context.Context, id uint64) (*model.AuditLog, error)
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.AuditLog, int64, error)
	FindByUserID(ctx context.Context, userID uint64) ([]*model.AuditLog, error)
	Count(ctx context.Context, params map[string]interface{}) (int64, error)
}

type auditLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepositoryImpl{db: db}
}

func (r *auditLogRepositoryImpl) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.AuditLog, error) {
	var log model.AuditLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	return &log, err
}

func (r *auditLogRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.AuditLog, int64, error) {
	var logs []*model.AuditLog
	db := r.db.WithContext(ctx)

	// 应用过滤条件
	if userID, ok := params["user_id"].(uint64); ok && userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if username, ok := params["username"].(string); ok && username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if role, ok := params["role"].(string); ok && role != "" {
		db = db.Where("role = ?", role)
	}
	if action, ok := params["action"].(string); ok && action != "" {
		db = db.Where("action = ?", action)
	}
	if resourceType, ok := params["resource_type"].(string); ok && resourceType != "" {
		db = db.Where("resource_type = ?", resourceType)
	}
	if resourceID, ok := params["resource_id"].(uint64); ok && resourceID > 0 {
		db = db.Where("resource_id = ?", resourceID)
	}
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		db = db.Where("created_at <= ?", endDate)
	}

	// 分页
	var count int64
	if err := db.Model(&model.AuditLog{}).Count(&count).Error; err != nil {
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

	err := db.Find(&logs).Error
	return logs, count, err
}

func (r *auditLogRepositoryImpl) FindByUserID(ctx context.Context, userID uint64) ([]*model.AuditLog, error) {
	var logs []*model.AuditLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepositoryImpl) Count(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// 应用过滤条件
	if userID, ok := params["user_id"].(uint64); ok && userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if username, ok := params["username"].(string); ok && username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if role, ok := params["role"].(string); ok && role != "" {
		db = db.Where("role = ?", role)
	}
	if action, ok := params["action"].(string); ok && action != "" {
		db = db.Where("action = ?", action)
	}
	if resourceType, ok := params["resource_type"].(string); ok && resourceType != "" {
		db = db.Where("resource_type = ?", resourceType)
	}
	if resourceID, ok := params["resource_id"].(uint64); ok && resourceID > 0 {
		db = db.Where("resource_id = ?", resourceID)
	}
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		db = db.Where("created_at <= ?", endDate)
	}

	return count, db.Count(&count).Error
}

// growth-partner/backend/internal/repository/child_repo.go

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type ChildRepository interface {
	Create(ctx context.Context, child *model.Child) error
	FindByUserID(ctx context.Context, userID uint64) (*model.Child, error)
	FindByID(ctx context.Context, id uint64) (*model.Child, error)
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.Child, int64, error)
	Update(ctx context.Context, child *model.Child) error
	BatchCreate(ctx context.Context, children []*model.Child) error
	UpdateStats(ctx context.Context, id uint64, growthDelta int) error
}

type childRepositoryImpl struct {
	db *gorm.DB
}

func NewChildRepository(db *gorm.DB) ChildRepository {
	return &childRepositoryImpl{db: db}
}

func (r *childRepositoryImpl) Create(ctx context.Context, child *model.Child) error {
	return r.db.WithContext(ctx).Create(child).Error
}

func (r *childRepositoryImpl) FindByUserID(ctx context.Context, userID uint64) (*model.Child, error) {
	var child model.Child
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&child).Error
	return &child, err
}

func (r *childRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.Child, error) {
	var child model.Child
	err := r.db.WithContext(ctx).First(&child, id).Error
	return &child, err
}

func (r *childRepositoryImpl) UpdateStats(ctx context.Context, id uint64, growthDelta int) error {
	return r.db.WithContext(ctx).Model(&model.Child{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_growth_points": gorm.Expr("current_growth_points + ?", growthDelta),
			"total_growth_points":   gorm.Expr("total_growth_points + ?", growthDelta),
			"total_behavior_count":  gorm.Expr("total_behavior_count + 1"),
		}).Error
}

// FindAll 查询所有学生（支持过滤、分页）
func (r *childRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.Child, int64, error) {
	var children []*model.Child
	db := r.db.WithContext(ctx)

	// 应用过滤条件
	if classID, ok := params["class_id"].(uint64); ok && classID > 0 {
		db = db.Where("class_id = ?", classID)
	}
	if studentID, ok := params["student_id"].(string); ok && studentID != "" {
		db = db.Where("student_no_enc = ?", studentID)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	// 分页
	var count int64
	if err := db.Model(&model.Child{}).Count(&count).Error; err != nil {
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

	err := db.Find(&children).Error
	return children, count, err
}

// Update 更新学生信息
func (r *childRepositoryImpl) Update(ctx context.Context, child *model.Child) error {
	return r.db.WithContext(ctx).Save(child).Error
}

// BatchCreate 批量创建学生
func (r *childRepositoryImpl) BatchCreate(ctx context.Context, children []*model.Child) error {
	return r.db.WithContext(ctx).CreateInBatches(children, 100).Error
}

// growth-partner/backend/internal/repository/school_repo.go
// 学校仓库接口和实现

package repository

import (
	"context"
	"growth-partner/internal/model"

	"gorm.io/gorm"
)

type SchoolRepository interface {
	Create(ctx context.Context, school *model.School) error
	FindByID(ctx context.Context, id uint64) (*model.School, error)
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.School, int64, error)
	Update(ctx context.Context, school *model.School) error
	UpdateStatus(ctx context.Context, id uint64, isActive bool) error
	Count(ctx context.Context, params map[string]interface{}) (int64, error)
}

type schoolRepositoryImpl struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolRepositoryImpl{db: db}
}

func (r *schoolRepositoryImpl) Create(ctx context.Context, school *model.School) error {
	return r.db.WithContext(ctx).Create(school).Error
}

func (r *schoolRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.School, error) {
	var school model.School
	err := r.db.WithContext(ctx).First(&school, id).Error
	return &school, err
}

func (r *schoolRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.School, int64, error) {
	var schools []*model.School
	db := r.db.WithContext(ctx)

	// 应用过滤条件
	if name, ok := params["name"].(string); ok && name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if district, ok := params["district"].(string); ok && district != "" {
		db = db.Where("district = ?", district)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	// 分页
	var count int64
	if err := db.Model(&model.School{}).Count(&count).Error; err != nil {
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

	err := db.Find(&schools).Error
	return schools, count, err
}

func (r *schoolRepositoryImpl) Update(ctx context.Context, school *model.School) error {
	return r.db.WithContext(ctx).Save(school).Error
}

func (r *schoolRepositoryImpl) UpdateStatus(ctx context.Context, id uint64, isActive bool) error {
	return r.db.WithContext(ctx).Model(&model.School{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *schoolRepositoryImpl) Count(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.School{})

	// 应用过滤条件
	if name, ok := params["name"].(string); ok && name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if district, ok := params["district"].(string); ok && district != "" {
		db = db.Where("district = ?", district)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	return count, db.Count(&count).Error
}

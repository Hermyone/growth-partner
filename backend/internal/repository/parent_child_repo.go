// growth-partner/backend/internal/repository/parent_child_repo.go
// 家长-学生绑定关系仓库接口和实现

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type ParentChildRepository interface {
	Create(ctx context.Context, binding *model.ParentChildRelation) error
	FindByID(ctx context.Context, id uint64) (*model.ParentChildRelation, error)
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.ParentChildRelation, int64, error)
	FindByParentID(ctx context.Context, parentID uint64) ([]*model.ParentChildRelation, error)
	FindByChildID(ctx context.Context, childID uint64) ([]*model.ParentChildRelation, error)
	Delete(ctx context.Context, id uint64) error
	SetAllNonPrimary(ctx context.Context, childID uint64) error
	Count(ctx context.Context, params map[string]interface{}) (int64, error)
}

type parentChildRepositoryImpl struct {
	db *gorm.DB
}

func NewParentChildRepository(db *gorm.DB) ParentChildRepository {
	return &parentChildRepositoryImpl{db: db}
}

func (r *parentChildRepositoryImpl) Create(ctx context.Context, binding *model.ParentChildRelation) error {
	return r.db.WithContext(ctx).Create(binding).Error
}

func (r *parentChildRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.ParentChildRelation, error) {
	var binding model.ParentChildRelation
	err := r.db.WithContext(ctx).First(&binding, id).Error
	return &binding, err
}

func (r *parentChildRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.ParentChildRelation, int64, error) {
	var bindings []*model.ParentChildRelation
	db := r.db.WithContext(ctx)

	// 应用过滤条件
	if parentUserID, ok := params["parent_user_id"].(uint64); ok && parentUserID > 0 {
		db = db.Where("parent_user_id = ?", parentUserID)
	}
	if childID, ok := params["child_id"].(uint64); ok && childID > 0 {
		db = db.Where("child_id = ?", childID)
	}
	if relationship, ok := params["relationship"].(string); ok && relationship != "" {
		db = db.Where("relationship = ?", relationship)
	}
	if isPrimary, ok := params["is_primary"].(bool); ok {
		db = db.Where("is_primary = ?", isPrimary)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	// 分页
	var count int64
	if err := db.Model(&model.ParentChildRelation{}).Count(&count).Error; err != nil {
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

	err := db.Find(&bindings).Error
	return bindings, count, err
}

func (r *parentChildRepositoryImpl) FindByParentID(ctx context.Context, parentID uint64) ([]*model.ParentChildRelation, error) {
	var bindings []*model.ParentChildRelation
	err := r.db.WithContext(ctx).Where("parent_user_id = ? AND is_active = ?", parentID, true).Find(&bindings).Error
	return bindings, err
}

func (r *parentChildRepositoryImpl) FindByChildID(ctx context.Context, childID uint64) ([]*model.ParentChildRelation, error) {
	var bindings []*model.ParentChildRelation
	err := r.db.WithContext(ctx).Where("child_id = ? AND is_active = ?", childID, true).Find(&bindings).Error
	return bindings, err
}

func (r *parentChildRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.ParentChildRelation{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *parentChildRepositoryImpl) Count(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.ParentChildRelation{})

	// 应用过滤条件
	if parentUserID, ok := params["parent_user_id"].(uint64); ok && parentUserID > 0 {
		db = db.Where("parent_user_id = ?", parentUserID)
	}
	if childID, ok := params["child_id"].(uint64); ok && childID > 0 {
		db = db.Where("child_id = ?", childID)
	}
	if relationship, ok := params["relationship"].(string); ok && relationship != "" {
		db = db.Where("relationship = ?", relationship)
	}
	if isPrimary, ok := params["is_primary"].(bool); ok {
		db = db.Where("is_primary = ?", isPrimary)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	return count, db.Count(&count).Error
}

// SetAllNonPrimary 将指定学生的所有绑定设置为非主绑定
func (r *parentChildRepositoryImpl) SetAllNonPrimary(ctx context.Context, childID uint64) error {
	return r.db.WithContext(ctx).Model(&model.ParentChildRelation{}).Where("child_id = ?", childID).Update("is_primary", false).Error
}

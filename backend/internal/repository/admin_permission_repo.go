// growth-partner/backend/internal/repository/admin_permission_repo.go
// 管理员权限仓库接口和实现

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type AdminPermissionRepository interface {
	Create(ctx context.Context, permission *model.AdminPermission) error
	FindByID(ctx context.Context, id uint64) (*model.AdminPermission, error)
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.AdminPermission, int64, error)
	FindByTeacherID(ctx context.Context, teacherID uint64) ([]*model.AdminPermission, error)
	Delete(ctx context.Context, id uint64) error
	BatchCreate(ctx context.Context, permissions []*model.AdminPermission) error
	Count(ctx context.Context, params map[string]interface{}) (int64, error)
}

type adminPermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminPermissionRepository(db *gorm.DB) AdminPermissionRepository {
	return &adminPermissionRepositoryImpl{db: db}
}

func (r *adminPermissionRepositoryImpl) Create(ctx context.Context, permission *model.AdminPermission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *adminPermissionRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.AdminPermission, error) {
	var permission model.AdminPermission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	return &permission, err
}

func (r *adminPermissionRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.AdminPermission, int64, error) {
	var permissions []*model.AdminPermission
	db := r.db.WithContext(ctx)

	// 应用过滤条件
	if teacherUserID, ok := params["teacher_user_id"].(uint64); ok && teacherUserID > 0 {
		db = db.Where("teacher_user_id = ?", teacherUserID)
	}
	if classID, ok := params["class_id"].(uint64); ok && classID > 0 {
		db = db.Where("class_id = ?", classID)
	}
	if permissionType, ok := params["permission_type"].(string); ok && permissionType != "" {
		db = db.Where("permission_type = ?", permissionType)
	}
	if schoolYear, ok := params["school_year"].(string); ok && schoolYear != "" {
		db = db.Where("school_year = ?", schoolYear)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	// 分页
	var count int64
	if err := db.Model(&model.AdminPermission{}).Count(&count).Error; err != nil {
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

	err := db.Find(&permissions).Error
	return permissions, count, err
}

func (r *adminPermissionRepositoryImpl) FindByTeacherID(ctx context.Context, teacherID uint64) ([]*model.AdminPermission, error) {
	var permissions []*model.AdminPermission
	err := r.db.WithContext(ctx).Where("teacher_user_id = ? AND is_active = ?", teacherID, true).Find(&permissions).Error
	return permissions, err
}

func (r *adminPermissionRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.AdminPermission{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *adminPermissionRepositoryImpl) BatchCreate(ctx context.Context, permissions []*model.AdminPermission) error {
	return r.db.WithContext(ctx).Create(&permissions).Error
}

func (r *adminPermissionRepositoryImpl) Count(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.AdminPermission{})

	// 应用过滤条件
	if teacherUserID, ok := params["teacher_user_id"].(uint64); ok && teacherUserID > 0 {
		db = db.Where("teacher_user_id = ?", teacherUserID)
	}
	if classID, ok := params["class_id"].(uint64); ok && classID > 0 {
		db = db.Where("class_id = ?", classID)
	}
	if permissionType, ok := params["permission_type"].(string); ok && permissionType != "" {
		db = db.Where("permission_type = ?", permissionType)
	}
	if schoolYear, ok := params["school_year"].(string); ok && schoolYear != "" {
		db = db.Where("school_year = ?", schoolYear)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	return count, db.Count(&count).Error
}

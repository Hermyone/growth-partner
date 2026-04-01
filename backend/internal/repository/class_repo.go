// growth-partner/backend/internal/repository/class_repo.go
// 班级仓储层：处理班级数据的持久化操作

package repository

import (
	"context"

	"growth-partner/internal/model"

	"gorm.io/gorm"
)

// ClassRepository 班级数据访问接口
type ClassRepository interface {
	Create(ctx context.Context, class *model.Class) error
	FindByID(ctx context.Context, id uint64) (*model.Class, error)
	FindByTeacherID(ctx context.Context, teacherID uint64) ([]*model.Class, error)
}

type classRepositoryImpl struct {
	db *gorm.DB
}

// NewClassRepository 创建班级仓储实例
func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepositoryImpl{db: db}
}

// Create 插入一个新的班级记录
func (r *classRepositoryImpl) Create(ctx context.Context, class *model.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

// FindByID 根据 ID 查询单个班级
func (r *classRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.Class, error) {
	var class model.Class
	err := r.db.WithContext(ctx).First(&class, id).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// FindByTeacherID 查询指定老师管理的所有活跃班级
func (r *classRepositoryImpl) FindByTeacherID(ctx context.Context, teacherID uint64) ([]*model.Class, error) {
	var classes []*model.Class
	err := r.db.WithContext(ctx).
		Where("teacher_id = ? AND is_active = ?", teacherID, true).
		Find(&classes).Error
	return classes, err
}

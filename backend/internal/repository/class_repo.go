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
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.Class, int64, error)
	Update(ctx context.Context, class *model.Class) error
	UpdateStatus(ctx context.Context, id uint64, isActive bool) error
	Count(ctx context.Context, params map[string]interface{}) (int64, error)
	FindByClassCode(ctx context.Context, classCode string) (*model.Class, error)
	BeginTx(ctx context.Context) *gorm.DB
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
		Where("homeroom_teacher_id = ? AND is_active = ?", teacherID, true).
		Find(&classes).Error
	return classes, err
}

// FindAll 查询所有班级（支持过滤、分页）
func (r *classRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.Class, int64, error) {
	var classes []*model.Class
	db := r.db.WithContext(ctx)

	// 应用过滤条件
	if schoolID, ok := params["school_id"].(uint64); ok && schoolID > 0 {
		db = db.Where("school_id = ?", schoolID)
	}
	if schoolYear, ok := params["school_year"].(string); ok && schoolYear != "" {
		db = db.Where("school_year = ?", schoolYear)
	}
	if grade, ok := params["grade"].(int); ok && grade > 0 {
		db = db.Where("grade = ?", grade)
	}
	if className, ok := params["class_name"].(string); ok && className != "" {
		db = db.Where("class_name LIKE ?", "%"+className+"%")
	}
	if classCode, ok := params["class_code"].(string); ok && classCode != "" {
		db = db.Where("class_code = ?", classCode)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	// 分页
	var count int64
	if err := db.Model(&model.Class{}).Count(&count).Error; err != nil {
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

	err := db.Find(&classes).Error
	return classes, count, err
}

// Update 更新班级信息
func (r *classRepositoryImpl) Update(ctx context.Context, class *model.Class) error {
	return r.db.WithContext(ctx).Save(class).Error
}

// UpdateStatus 更新班级状态
func (r *classRepositoryImpl) UpdateStatus(ctx context.Context, id uint64, isActive bool) error {
	return r.db.WithContext(ctx).Model(&model.Class{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// Count 统计班级数量
func (r *classRepositoryImpl) Count(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.Class{})

	// 应用过滤条件
	if schoolID, ok := params["school_id"].(uint64); ok && schoolID > 0 {
		db = db.Where("school_id = ?", schoolID)
	}
	if schoolYear, ok := params["school_year"].(string); ok && schoolYear != "" {
		db = db.Where("school_year = ?", schoolYear)
	}
	if grade, ok := params["grade"].(int); ok && grade > 0 {
		db = db.Where("grade = ?", grade)
	}
	if className, ok := params["class_name"].(string); ok && className != "" {
		db = db.Where("class_name LIKE ?", "%"+className+"%")
	}
	if classCode, ok := params["class_code"].(string); ok && classCode != "" {
		db = db.Where("class_code = ?", classCode)
	}
	if isActive, ok := params["is_active"].(bool); ok {
		db = db.Where("is_active = ?", isActive)
	}

	return count, db.Count(&count).Error
}

// FindByClassCode 根据班级代码查询班级
func (r *classRepositoryImpl) FindByClassCode(ctx context.Context, classCode string) (*model.Class, error) {
	var class model.Class
	err := r.db.WithContext(ctx).Where("class_code = ?", classCode).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// BeginTx 开始一个新的事务
func (r *classRepositoryImpl) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

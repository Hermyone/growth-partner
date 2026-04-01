// growth-partner/backend/internal/repository/behavior_repo.go
// 行为记录数据访问层：封装所有的数据库操作

package repository

import (
	"context"

	"growth-partner/internal/model"

	"gorm.io/gorm"
)

// BehaviorRepository 行为记录仓库接口
type BehaviorRepository interface {
	Create(ctx context.Context, record *model.BehaviorRecord) error
	FindByChildID(ctx context.Context, childID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error)
	FindByClassID(ctx context.Context, classID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error)
}

type behaviorRepositoryImpl struct {
	db *gorm.DB
}

// NewBehaviorRepository 创建行为记录仓库实例
func NewBehaviorRepository(db *gorm.DB) BehaviorRepository {
	return &behaviorRepositoryImpl{db: db}
}

// Create 插入一条正向行为记录
func (r *behaviorRepositoryImpl) Create(ctx context.Context, record *model.BehaviorRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// FindByChildID 分页获取某个学生的所有行为记录（按时间倒序）
func (r *behaviorRepositoryImpl) FindByChildID(ctx context.Context, childID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error) {
	var records []*model.BehaviorRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BehaviorRecord{}).Where("child_id = ?", childID)

	// 先查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查分页数据
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// FindByClassID 分页获取某个班级的所有行为记录（班主任用，按时间倒序）
func (r *behaviorRepositoryImpl) FindByClassID(ctx context.Context, classID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error) {
	var records []*model.BehaviorRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BehaviorRecord{}).Where("class_id = ?", classID)

	// 先查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查分页数据
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

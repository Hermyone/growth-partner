// growth-partner/backend/internal/repository/growth_repo.go

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type GrowthRepository interface {
	Create(ctx context.Context, record *model.GrowthRecord) error
	FindByChildID(ctx context.Context, childID uint64, limit, offset int) ([]*model.GrowthRecord, int64, error)
}

type growthRepositoryImpl struct {
	db *gorm.DB
}

func NewGrowthRepository(db *gorm.DB) GrowthRepository {
	return &growthRepositoryImpl{db: db}
}

func (r *growthRepositoryImpl) Create(ctx context.Context, record *model.GrowthRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *growthRepositoryImpl) FindByChildID(ctx context.Context, childID uint64, limit, offset int) ([]*model.GrowthRecord, int64, error) {
	var list []*model.GrowthRecord
	var total int64
	db := r.db.WithContext(ctx).Model(&model.GrowthRecord{}).Where("child_id = ?", childID)
	db.Count(&total)
	err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

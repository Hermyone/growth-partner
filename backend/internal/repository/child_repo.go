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

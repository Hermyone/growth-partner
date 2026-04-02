// growth-partner/backend/internal/repository/broadcast_repo.go
// 广播仓库接口和实现

package repository

import (
	"context"
	"growth-partner/internal/model"

	"gorm.io/gorm"
)

type BroadcastRepository interface {
	Create(ctx context.Context, broadcast *model.Broadcast) error
	FindByID(ctx context.Context, id uint64) (*model.Broadcast, error)
	FindByCreatorID(ctx context.Context, creatorID uint64) ([]*model.Broadcast, error)
	Update(ctx context.Context, broadcast *model.Broadcast) error
	Delete(ctx context.Context, id uint64) error
}

type broadcastRepositoryImpl struct {
	db *gorm.DB
}

func NewBroadcastRepository(db *gorm.DB) BroadcastRepository {
	return &broadcastRepositoryImpl{db: db}
}

func (r *broadcastRepositoryImpl) Create(ctx context.Context, broadcast *model.Broadcast) error {
	return r.db.WithContext(ctx).Create(broadcast).Error
}

func (r *broadcastRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.Broadcast, error) {
	var broadcast model.Broadcast
	err := r.db.WithContext(ctx).First(&broadcast, id).Error
	return &broadcast, err
}

func (r *broadcastRepositoryImpl) FindByCreatorID(ctx context.Context, creatorID uint64) ([]*model.Broadcast, error) {
	var broadcasts []*model.Broadcast
	err := r.db.WithContext(ctx).Where("created_by = ?", creatorID).Order("created_at DESC").Find(&broadcasts).Error
	return broadcasts, err
}

func (r *broadcastRepositoryImpl) Update(ctx context.Context, broadcast *model.Broadcast) error {
	return r.db.WithContext(ctx).Save(broadcast).Error
}

func (r *broadcastRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Broadcast{}, id).Error
}

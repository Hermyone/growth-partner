// growth-partner/backend/internal/repository/battle_repo.go

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type BattleRepository interface {
	Create(ctx context.Context, record *model.BattleRecord) error
	FindByChildID(ctx context.Context, childID uint64, limit, offset int) ([]*model.BattleRecord, int64, error)
}

type battleRepositoryImpl struct {
	db *gorm.DB
}

func NewBattleRepository(db *gorm.DB) BattleRepository {
	return &battleRepositoryImpl{db: db}
}

func (r *battleRepositoryImpl) Create(ctx context.Context, record *model.BattleRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *battleRepositoryImpl) FindByChildID(ctx context.Context, childID uint64, limit, offset int) ([]*model.BattleRecord, int64, error) {
	var list []*model.BattleRecord
	var total int64
	db := r.db.WithContext(ctx).Model(&model.BattleRecord{}).
		Where("player_a_child_id = ? OR player_b_child_id = ?", childID, childID)
	db.Count(&total)
	err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

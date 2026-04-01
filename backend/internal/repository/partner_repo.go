// growth-partner/backend/internal/repository/partner_repo.go

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type PartnerRepository interface {
	Create(ctx context.Context, partner *model.Partner) error
	FindByChildID(ctx context.Context, childID uint64) (*model.Partner, error)
	UpdateNickname(ctx context.Context, childID uint64, nickname string) error
	UpdateGrowthState(ctx context.Context, partner *model.Partner) error
}

type partnerRepositoryImpl struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) PartnerRepository {
	return &partnerRepositoryImpl{db: db}
}

func (r *partnerRepositoryImpl) Create(ctx context.Context, partner *model.Partner) error {
	return r.db.WithContext(ctx).Create(partner).Error
}

func (r *partnerRepositoryImpl) FindByChildID(ctx context.Context, childID uint64) (*model.Partner, error) {
	var partner model.Partner
	err := r.db.WithContext(ctx).Where("child_id = ?", childID).First(&partner).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &partner, err
}

func (r *partnerRepositoryImpl) UpdateNickname(ctx context.Context, childID uint64, nickname string) error {
	return r.db.WithContext(ctx).Model(&model.Partner{}).Where("child_id = ?", childID).Update("nickname", nickname).Error
}

func (r *partnerRepositoryImpl) UpdateGrowthState(ctx context.Context, p *model.Partner) error {
	return r.db.WithContext(ctx).Model(p).Updates(map[string]interface{}{
		"growth_points":      p.GrowthPoints,
		"current_stage":      p.CurrentStage,
		"evolution_count":    p.EvolutionCount,
		"interaction_level":  p.InteractionLevel,
		"last_evolved_at":    p.LastEvolvedAt,
		"first_evolution_at": p.FirstEvolutionAt,
	}).Error
}

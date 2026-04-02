// growth-partner/backend/internal/repository/challenge_repo.go
// 集体挑战仓库接口和实现

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type ChallengeRepository interface {
	Create(ctx context.Context, challenge *model.Challenge) error
	FindByID(ctx context.Context, id uint64) (*model.Challenge, error)
	FindByClassID(ctx context.Context, classID uint64) ([]*model.Challenge, error)
	Update(ctx context.Context, challenge *model.Challenge) error
	Complete(ctx context.Context, id uint64, completedBy uint64) error
	Delete(ctx context.Context, id uint64) error
}

type challengeRepositoryImpl struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &challengeRepositoryImpl{db: db}
}

func (r *challengeRepositoryImpl) Create(ctx context.Context, challenge *model.Challenge) error {
	return r.db.WithContext(ctx).Create(challenge).Error
}

func (r *challengeRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.Challenge, error) {
	var challenge model.Challenge
	err := r.db.WithContext(ctx).First(&challenge, id).Error
	return &challenge, err
}

func (r *challengeRepositoryImpl) FindByClassID(ctx context.Context, classID uint64) ([]*model.Challenge, error) {
	var challenges []*model.Challenge
	err := r.db.WithContext(ctx).Where("class_id = ?", classID).Order("created_at DESC").Find(&challenges).Error
	return challenges, err
}

func (r *challengeRepositoryImpl) Update(ctx context.Context, challenge *model.Challenge) error {
	return r.db.WithContext(ctx).Save(challenge).Error
}

func (r *challengeRepositoryImpl) Complete(ctx context.Context, id uint64, completedBy uint64) error {
	return r.db.WithContext(ctx).Model(&model.Challenge{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       model.ChallengeStatusCompleted,
		"completed_by": completedBy,
	}).Error
}

func (r *challengeRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Challenge{}, id).Error
}

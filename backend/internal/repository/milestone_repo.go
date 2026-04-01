// growth-partner/backend/internal/repository/milestone_repo.go
// 里程碑仓储：记录孩子成长过程中的重大成就，用于生成成长年历

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type MilestoneRepository interface {
	Create(ctx context.Context, ms *model.Milestone) error
	FindByChildID(ctx context.Context, childID uint64) ([]*model.Milestone, error)
	MarkAsNotified(ctx context.Context, id uint64) error
}

type milestoneRepositoryImpl struct {
	db *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &milestoneRepositoryImpl{db: db}
}

// Create 记录一个新的里程碑成就
func (r *milestoneRepositoryImpl) Create(ctx context.Context, ms *model.Milestone) error {
	return r.db.WithContext(ctx).Create(ms).Error
}

// FindByChildID 获取孩子获得的所有里程碑（用于渲染成长画卷）
func (r *milestoneRepositoryImpl) FindByChildID(ctx context.Context, childID uint64) ([]*model.Milestone, error) {
	var list []*model.Milestone
	err := r.db.WithContext(ctx).
		Where("child_id = ?", childID).
		Order("created_at ASC").
		Find(&list).Error
	return list, err
}

// MarkAsNotified 标记为已在前端弹出过庆祝动画，防止重复弹出
func (r *milestoneRepositoryImpl) MarkAsNotified(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Milestone{}).
		Where("id = ?", id).
		Update("is_notified", true).Error
}

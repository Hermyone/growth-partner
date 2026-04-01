// growth-partner/backend/internal/service/blindbox_service.go
// 盲盒服务：管理奖池与抽取逻辑

package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
	"math/rand"
	"time"
)

type BlindboxService interface {
	GetPool(ctx context.Context, classID uint64) ([]*model.BlindBoxPool, error)
	Draw(ctx context.Context, childID, classID uint64) (*model.BlindBoxDraw, *model.BlindBoxPool, error)
}

type blindboxServiceImpl struct {
	db *gorm.DB
}

func NewBlindboxService(repo repository.BlindboxRepository, db *gorm.DB) BlindboxService {
	return &blindboxServiceImpl{db: db}
}

func (s *blindboxServiceImpl) GetPool(ctx context.Context, classID uint64) ([]*model.BlindBoxPool, error) {
	var pool []*model.BlindBoxPool
	err := s.db.WithContext(ctx).Where("class_id = ? AND is_active = ? AND stock != 0", classID, true).Find(&pool).Error
	return pool, err
}

func (s *blindboxServiceImpl) Draw(ctx context.Context, childID, classID uint64) (*model.BlindBoxDraw, *model.BlindBoxPool, error) {
	pool, err := s.GetPool(ctx, classID)
	if err != nil || len(pool) == 0 {
		return nil, nil, fmt.Errorf("奖池为空")
	}

	// 随机抽取
	rand.Seed(time.Now().UnixNano())
	reward := pool[rand.Intn(len(pool))]

	draw := &model.BlindBoxDraw{
		ChildID: childID,
		ClassID: classID,
		PoolID:  reward.ID,
		DrawnAt: time.Now(),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(draw).Error; err != nil {
			return err
		}
		// 扣减库存
		if reward.Stock > 0 {
			return tx.Model(reward).Update("stock", reward.Stock-1).Error
		}
		return nil
	})

	return draw, reward, err
}

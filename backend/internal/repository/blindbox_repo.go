// growth-partner/backend/internal/repository/blindbox_repo.go
// 盲盒仓储层：管理班级奖励池和抽取记录

package repository

import (
	"context"
	"fmt"

	"growth-partner/internal/model"

	"gorm.io/gorm"
)

type BlindboxRepository interface {
	// 奖池管理
	GetPoolByClassID(ctx context.Context, classID uint64) ([]*model.BlindBoxPool, error)
	CreatePoolEntry(ctx context.Context, pool *model.BlindBoxPool) error
	DeletePoolEntry(ctx context.Context, id uint64) error
	FindPoolByID(ctx context.Context, id uint64) (*model.BlindBoxPool, error)

	// 抽取逻辑
	CreateDraw(ctx context.Context, draw *model.BlindBoxDraw) error
	// 事务：记录抽取并更新库存
	ExecuteDrawTransaction(ctx context.Context, draw *model.BlindBoxDraw, poolID uint64) error
}

type blindboxRepositoryImpl struct {
	db *gorm.DB
}

func NewBlindboxRepository(db *gorm.DB) BlindboxRepository {
	return &blindboxRepositoryImpl{db: db}
}

// GetPoolByClassID 获取班级当前可用的盲盒奖励（有库存且活跃）
func (r *blindboxRepositoryImpl) GetPoolByClassID(ctx context.Context, classID uint64) ([]*model.BlindBoxPool, error) {
	var pool []*model.BlindBoxPool
	err := r.db.WithContext(ctx).
		Where("class_id = ? AND is_active = ? AND (stock > 0 OR stock = -1)", classID, true).
		Find(&pool).Error
	return pool, err
}

func (r *blindboxRepositoryImpl) CreatePoolEntry(ctx context.Context, pool *model.BlindBoxPool) error {
	return r.db.WithContext(ctx).Create(pool).Error
}

func (r *blindboxRepositoryImpl) DeletePoolEntry(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BlindBoxPool{}, id).Error
}

func (r *blindboxRepositoryImpl) FindPoolByID(ctx context.Context, id uint64) (*model.BlindBoxPool, error) {
	var pool model.BlindBoxPool
	err := r.db.WithContext(ctx).First(&pool, id).Error
	return &pool, err
}

func (r *blindboxRepositoryImpl) CreateDraw(ctx context.Context, draw *model.BlindBoxDraw) error {
	return r.db.WithContext(ctx).Create(draw).Error
}

// ExecuteDrawTransaction 核心事务：保存抽取记录并原子化扣减库存
func (r *blindboxRepositoryImpl) ExecuteDrawTransaction(ctx context.Context, draw *model.BlindBoxDraw, poolID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建抽取记录
		if err := tx.Create(draw).Error; err != nil {
			return err
		}

		// 2. 锁定并查询奖项
		var pool model.BlindBoxPool
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&pool, poolID).Error; err != nil {
			return err
		}

		// 3. 检查库存
		if pool.Stock == 0 {
			return fmt.Errorf("该奖励已被抽完")
		}

		// 4. 更新库存（-1 表示无限，不扣减）
		if pool.Stock > 0 {
			if err := tx.Model(&pool).Update("stock", gorm.Expr("stock - 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

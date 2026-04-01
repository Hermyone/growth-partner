// growth-partner/backend/internal/service/battle_service.go
// 对战服务：处理知识对战逻辑与房间管理

package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

type BattleService interface {
	CreateRoom(ctx context.Context, childID uint64, subject model.BattleSubject) (string, error)
	GetBattleHistory(ctx context.Context, childID uint64, limit, offset int) ([]*model.BattleRecord, int64, error)
}

type battleServiceImpl struct {
	battleRepo repository.BattleRepository
	rdb        *redis.Client
}

func NewBattleService(repo repository.BattleRepository, rdb *redis.Client) BattleService {
	return &battleServiceImpl{battleRepo: repo, rdb: rdb}
}

func (s *battleServiceImpl) CreateRoom(ctx context.Context, childID uint64, subject model.BattleSubject) (string, error) {
	roomID := fmt.Sprintf("room:%s:%d", subject, childID)
	// 在 Redis 中注册房间信息，后续 WebSocket 会根据此 ID 加入
	err := s.rdb.Set(ctx, roomID, "waiting", 0).Err()
	return roomID, err
}

func (s *battleServiceImpl) GetBattleHistory(ctx context.Context, childID uint64, limit, offset int) ([]*model.BattleRecord, int64, error) {
	return s.battleRepo.FindByChildID(ctx, childID, limit, offset)
}

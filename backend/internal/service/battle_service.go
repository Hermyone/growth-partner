// growth-partner/backend/internal/service/battle_service.go
// 知识对战服务层

package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"growth-partner/config"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// BattleService 对战服务接口
type BattleService interface {
	// 获取可用对战科目列表
	GetBattleSubjects(ctx context.Context) (map[string]string, error)

	// 创建对战房间
	CreateBattleRoom(ctx context.Context, studentID uint64, subject string) (*model.BattleRoom, error)

	// 加入对战房间
	JoinBattleRoom(ctx context.Context, studentID uint64, roomCode string) (*model.BattleRoom, error)

	// 查询房间状态
	GetBattleRoom(ctx context.Context, roomCode string) (*model.BattleRoom, error)

	// 查看对战历史
	GetBattleHistory(ctx context.Context, studentID uint64) ([]*model.BattleRecord, error)

	// 对战复盘
	GetBattleReview(ctx context.Context, studentID uint64, roomID string) (map[string]interface{}, error)
}

// battleServiceImpl 对战服务实现
type battleServiceImpl struct {
	battleRepo   repository.BattleRepository
	questionRepo repository.QuestionRepository
	childRepo    repository.ChildRepository
}

// NewBattleService 创建对战服务实例
func NewBattleService(
	battleRepo repository.BattleRepository,
	questionRepo repository.QuestionRepository,
	childRepo repository.ChildRepository,
) BattleService {
	return &battleServiceImpl{
		battleRepo:   battleRepo,
		questionRepo: questionRepo,
		childRepo:    childRepo,
	}
}

// ─── 对战模块核心功能 ────────────────────────────────────────

func (s *battleServiceImpl) GetBattleSubjects(ctx context.Context) (map[string]string, error) {
	// 获取可用对战科目列表
	// 暂时返回模拟数据
	return map[string]string{
		"poetry":  "诗词",
		"math":    "速算",
		"general": "百科",
		"english": "英语",
		"custom":  "自定义题库",
	}, nil
}

func (s *battleServiceImpl) CreateBattleRoom(ctx context.Context, studentID uint64, subject string) (*model.BattleRoom, error) {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("学生不存在")
	}

	// 生成6位邀请码
	roomCode := generateRoomCode()

	// 创建对战房间
	room := &model.BattleRoom{
		RoomCode:  roomCode,
		Subject:   model.BattleSubject(subject),
		Status:    model.BattleRoomStatusWaiting,
		CreatorID: studentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 暂时返回房间，实际应该保存到数据库
	return room, nil
}

func (s *battleServiceImpl) JoinBattleRoom(ctx context.Context, studentID uint64, roomCode string) (*model.BattleRoom, error) {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("学生不存在")
	}

	// 暂时返回房间，实际应该从数据库查找
	room := &model.BattleRoom{
		RoomCode:      roomCode,
		Subject:       model.BattleSubject("general"),
		Status:        model.BattleRoomStatusBattle,
		CreatorID:     config.Get().Battle.DefaultCreatorID, // 从配置中获取
		ParticipantID: studentID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return room, nil
}

func (s *battleServiceImpl) GetBattleRoom(ctx context.Context, roomCode string) (*model.BattleRoom, error) {
	// 暂时返回房间，实际应该从数据库查找
	room := &model.BattleRoom{
		RoomCode:   roomCode,
		Subject:    model.BattleSubject("general"),
		Status:     model.BattleRoomStatusWaiting,
		CreatorID:  config.Get().Battle.DefaultCreatorID, // 从配置中获取
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return room, nil
}

func (s *battleServiceImpl) GetBattleHistory(ctx context.Context, studentID uint64) ([]*model.BattleRecord, error) {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("学生不存在")
	}

	// 获取对战历史
	records, _, err := s.battleRepo.FindByChildID(ctx, studentID, 10, 0)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *battleServiceImpl) GetBattleReview(ctx context.Context, studentID uint64, roomID string) (map[string]interface{}, error) {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("学生不存在")
	}

	// 暂时返回模拟数据，实际应该从数据库查找
	review := map[string]interface{}{
		"room_id":         roomID,
		"subject":         "general",
		"start_time":      time.Now(),
		"end_time":        time.Now(),
		"player_a_id":     studentID,
		"player_b_id":     uint64(2),
		"player_a_score":  8,
		"player_b_score":  6,
		"winner_id":       &studentID,
		"honor_badge":     "胜利徽章",
		"duration":        60,
		"player_a_growth": 3,
		"player_b_growth": 3,
	}

	return review, nil
}

// generateRoomCode 生成随机邀请码
func generateRoomCode() string {
	const charset = "0123456789"
	codeLength := config.Get().Battle.RoomCodeLength
	code := make([]byte, codeLength)
	rand.Seed(time.Now().UnixNano())
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

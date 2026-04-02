// growth-partner/backend/internal/service/behavior_service.go
// 行为记录业务逻辑：处理正向行为的录入，并联动伙伴系统增加成长值

package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

var (
	ErrInvalidDimension = errors.New("无效的行为维度")
	ErrInvalidGrowth    = errors.New("成长值超出该维度的合理范围")
)

// BehaviorService 行为记录服务接口
type BehaviorService interface {
	// 记录学生正向行为（核心接口）
	RecordBehavior(ctx context.Context, req RecordBehaviorRequest) (*model.BehaviorRecord, *AddGrowthPointsResult, error)
	// 获取学生的行为记录
	GetChildBehaviors(ctx context.Context, childID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error)
	// 获取班级的行为记录
	GetClassBehaviors(ctx context.Context, classID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error)
}

// RecordBehaviorRequest 记录行为请求 DTO
type RecordBehaviorRequest struct {
	ChildID      uint64
	ClassID      uint64
	SchoolYear   string // 所属学年
	RecorderID   uint64 // 记录人ID (老师/家长)
	RecorderRole model.RecorderRole
	Dimension    model.BehaviorDimension
	Description  string
	GrowthValue  int // 获得的成长值
}

type behaviorServiceImpl struct {
	behaviorRepo repository.BehaviorRepository
	partnerSvc   PartnerService
	broadcastSvc BroadcastService
}

// NewBehaviorService 创建行为服务实例
func NewBehaviorService(
	repo repository.BehaviorRepository,
	partnerSvc PartnerService,
	broadcastSvc BroadcastService,
) BehaviorService {
	return &behaviorServiceImpl{
		behaviorRepo: repo,
		partnerSvc:   partnerSvc,
		broadcastSvc: broadcastSvc,
	}
}

// RecordBehavior 记录正向行为并联动伙伴成长
func (s *behaviorServiceImpl) RecordBehavior(ctx context.Context, req RecordBehaviorRequest) (*model.BehaviorRecord, *AddGrowthPointsResult, error) {
	// 1. 验证行为维度与分值范围
	allowedRange, ok := model.DimensionGrowthRange[req.Dimension]
	if !ok {
		return nil, nil, ErrInvalidDimension
	}
	if req.GrowthValue < allowedRange[0] || req.GrowthValue > allowedRange[1] {
		// 容错处理：如果不合法，默认给该维度最低分，体现“只记录进步”的原则
		log.Printf("[BehaviorService] 警告：分值 %d 异常，已重置为最低分 %d", req.GrowthValue, allowedRange[0])
		req.GrowthValue = allowedRange[0]
	}

	// 2. 生成伙伴的鼓励话语 (未来可接入预置回复库)
	partnerMessage := fmt.Sprintf("太棒啦！因为你的努力，我又成长了 %d 点！", req.GrowthValue)

	// 3. 创建行为记录实体
	record := &model.BehaviorRecord{
		ChildID:        req.ChildID,
		ClassID:        req.ClassID,
		SchoolYear:     req.SchoolYear,
		RecorderUserID: req.RecorderID,
		RecorderRole:   req.RecorderRole,
		Dimension:      req.Dimension,
		Description:    req.Description,
		GrowthValue:    req.GrowthValue,
		PartnerMessage: partnerMessage,
	}

	// 4. 保存记录到数据库
	if err := s.behaviorRepo.Create(ctx, record); err != nil {
		return nil, nil, fmt.Errorf("保存行为记录失败: %w", err)
	}

	// 5. 联动伙伴系统，增加成长值并检测进化
	growthReq := AddGrowthPointsRequest{
		ChildID:    req.ChildID,
		BehaviorID: record.ID,
		SourceType: model.GrowthSourceBehavior,
		Delta:      req.GrowthValue,
		Remark:     fmt.Sprintf("因为 [%s] 获得了进步", getDimensionName(req.Dimension)),
	}

	growthResult, err := s.partnerSvc.AddGrowthPoints(ctx, growthReq)
	if err != nil {
		log.Printf("[BehaviorService] 警告：加分失败，但不影响行为记录: %v", err)
	}

	// 6. 异步推送实时消息
	go func() {
		if s.broadcastSvc != nil {
			_ = s.broadcastSvc.BroadcastPartnerMessage(context.Background(), req.ChildID, partnerMessage)
		}
	}()

	return record, growthResult, nil
}

func (s *behaviorServiceImpl) GetChildBehaviors(ctx context.Context, childID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error) {
	return s.behaviorRepo.FindByChildID(ctx, childID, limit, offset)
}

func (s *behaviorServiceImpl) GetClassBehaviors(ctx context.Context, classID uint64, limit, offset int) ([]*model.BehaviorRecord, int64, error) {
	return s.behaviorRepo.FindByClassID(ctx, classID, limit, offset)
}

// translateDimension 辅助函数
func getDimensionName(dim model.BehaviorDimension) string {
	switch dim {
	case model.DimVirtue:
		return "德馨"
	case model.DimStudy:
		return "智睿"
	case model.DimSport:
		return "体健"
	case model.DimArt:
		return "美雅"
	case model.DimLabor:
		return "劳朴"
	default:
		return "未知"
	}
}

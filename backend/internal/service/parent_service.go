// growth-partner/backend/internal/service/parent_service.go
// 家长端服务层

package service

import (
	"context"
	"errors"

	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// ParentService 家长端服务接口
type ParentService interface {
	// 获取绑定的孩子列表
	GetMyChildren(ctx context.Context, parentID uint64) ([]*model.Child, error)
	
	// 查看孩子当前伙伴状态
	GetChildPartner(ctx context.Context, childID, parentID uint64) (*model.Partner, error)
	
	// 查看孩子的历史伙伴列表
	GetChildPartners(ctx context.Context, childID, parentID uint64) ([]*model.Partner, error)
	
	// 查看孩子的正向行为记录
	GetChildBehaviors(ctx context.Context, childID, parentID uint64, params map[string]interface{}) ([]*model.BehaviorRecord, int64, error)
	
	// 查看孩子收到的伙伴鼓励广播
	GetChildBroadcasts(ctx context.Context, childID, parentID uint64) ([]*model.Broadcast, error)
	
	// 查看孩子的里程碑贴纸
	GetChildMilestones(ctx context.Context, childID, parentID uint64) ([]*model.Milestone, error)
	
	// 查看孩子本月/历史月度成长卡
	GetChildMonthlyCard(ctx context.Context, childID, parentID uint64, month string) (map[string]interface{}, error)
	
	// 查看孩子年度成长画卷
	GetChildAnnualReport(ctx context.Context, childID, parentID uint64, year string) (map[string]interface{}, error)
	
	// 查看孩子的对战参与记录
	GetChildBattles(ctx context.Context, childID, parentID uint64) ([]*model.BattleRecord, error)
}

// parentServiceImpl 家长端服务实现
type parentServiceImpl struct {
	parentChildRepo repository.ParentChildRepository
	childRepo       repository.ChildRepository
	partnerRepo     repository.PartnerRepository
	behaviorRepo    repository.BehaviorRepository
	broadcastRepo   repository.BroadcastRepository
	milestoneRepo   repository.MilestoneRepository
	battleRepo      repository.BattleRepository
}

// NewParentService 创建家长端服务实例
func NewParentService(
	parentChildRepo repository.ParentChildRepository,
	childRepo repository.ChildRepository,
	partnerRepo repository.PartnerRepository,
	behaviorRepo repository.BehaviorRepository,
	broadcastRepo repository.BroadcastRepository,
	milestoneRepo repository.MilestoneRepository,
	battleRepo repository.BattleRepository,
) ParentService {
	return &parentServiceImpl{
		parentChildRepo: parentChildRepo,
		childRepo:       childRepo,
		partnerRepo:     partnerRepo,
		behaviorRepo:    behaviorRepo,
		broadcastRepo:   broadcastRepo,
		milestoneRepo:   milestoneRepo,
		battleRepo:      battleRepo,
	}
}

// ─── 家长端核心功能 ──────────────────────────────────────────

func (s *parentServiceImpl) GetMyChildren(ctx context.Context, parentID uint64) ([]*model.Child, error) {
	// 获取家长绑定的孩子列表
	bindings, err := s.parentChildRepo.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	
	// 遍历绑定关系，获取孩子详情
	children := make([]*model.Child, 0, len(bindings))
	for _, binding := range bindings {
		child, err := s.childRepo.FindByID(ctx, binding.ChildID)
		if err != nil {
			return nil, err
		}
		if child != nil {
			children = append(children, child)
		}
	}
	
	return children, nil
}

func (s *parentServiceImpl) GetChildPartner(ctx context.Context, childID, parentID uint64) (*model.Partner, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子的伙伴
	partner, err := s.partnerRepo.FindByChildID(ctx, childID)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return nil, errors.New("孩子还没有伙伴")
	}
	
	return partner, nil
}

func (s *parentServiceImpl) GetChildPartners(ctx context.Context, childID, parentID uint64) ([]*model.Partner, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子的伙伴（由于PartnerRepository只提供了FindByChildID方法，暂时只返回当前伙伴）
	partner, err := s.partnerRepo.FindByChildID(ctx, childID)
	if err != nil {
		return nil, err
	}
	
	if partner == nil {
		return []*model.Partner{}, nil
	}
	
	return []*model.Partner{partner}, nil
}

func (s *parentServiceImpl) GetChildBehaviors(ctx context.Context, childID, parentID uint64, params map[string]interface{}) ([]*model.BehaviorRecord, int64, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, 0, err
	}
	
	// 解析分页参数
	limit := 10
	offset := 0
	if l, ok := params["limit"].(int); ok && l > 0 {
		limit = l
	}
	if page, ok := params["page"].(int); ok && page > 0 {
		offset = (page - 1) * limit
	}
	
	// 获取孩子的行为记录
	return s.behaviorRepo.FindByChildID(ctx, childID, limit, offset)
}

func (s *parentServiceImpl) GetChildBroadcasts(ctx context.Context, childID, parentID uint64) ([]*model.Broadcast, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子收到的广播消息（暂时返回空数组，后续需要根据实际的广播记录模型来实现）
	return []*model.Broadcast{}, nil
}

func (s *parentServiceImpl) GetChildMilestones(ctx context.Context, childID, parentID uint64) ([]*model.Milestone, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子的里程碑列表（暂时返回空数组，后续需要根据实际的里程碑模型来实现）
	return []*model.Milestone{}, nil
}

func (s *parentServiceImpl) GetChildMonthlyCard(ctx context.Context, childID, parentID uint64, month string) (map[string]interface{}, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子的月度成长卡（暂时返回空map，后续需要根据实际的月度成长卡模型来实现）
	return map[string]interface{}{"month": month, "data": nil}, nil
}

func (s *parentServiceImpl) GetChildAnnualReport(ctx context.Context, childID, parentID uint64, year string) (map[string]interface{}, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子的年度成长画卷（暂时返回空map，后续需要根据实际的年度报告模型来实现）
	return map[string]interface{}{"year": year, "data": nil}, nil
}

func (s *parentServiceImpl) GetChildBattles(ctx context.Context, childID, parentID uint64) ([]*model.BattleRecord, error) {
	// 检查家长是否有权限查看该孩子
	if err := s.checkParentChildBinding(ctx, parentID, childID); err != nil {
		return nil, err
	}
	
	// 获取孩子的对战参与记录（暂时返回空数组，后续需要根据实际的对战记录模型来实现）
	return []*model.BattleRecord{}, nil
}

// checkParentChildBinding 检查家长与孩子的绑定关系
func (s *parentServiceImpl) checkParentChildBinding(ctx context.Context, parentID, childID uint64) error {
	// 获取家长与孩子的绑定关系
	bindings, err := s.parentChildRepo.FindByParentID(ctx, parentID)
	if err != nil {
		return err
	}
	
	// 检查是否存在绑定关系
	for _, binding := range bindings {
		if binding.ChildID == childID {
			return nil
		}
	}
	
	return errors.New("您没有权限查看该孩子的信息")
}

// growth-partner/backend/internal/service/sunshine_service.go
// 阳光章服务层

package service

import (
	"context"
	"errors"
	"time"

	"growth-partner/config"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// SunshineService 阳光章服务接口
type SunshineService interface {
	// 七色配置管理（管理员）
	GetSunshineColors(ctx context.Context, schoolID uint64) ([]*model.SunshineColor, error)
	CreateSunshineColor(ctx context.Context, color *model.SunshineColor) error
	UpdateSunshineColor(ctx context.Context, color *model.SunshineColor) error

	// 盖章操作
	StampSunshine(ctx context.Context, teacherID, studentID, colorID uint64, subject string) error
	GetClassStamps(ctx context.Context, teacherID, classID uint64, params map[string]interface{}) ([]*model.SunshineStamp, int64, error)
	GetStudentStamps(ctx context.Context, studentID uint64) (map[string]interface{}, error)

	// 之星评选
	EvaluateSunshineAwards(ctx context.Context, teacherID, classID uint64, period string) error
	GetSunshineAwards(ctx context.Context, teacherID, classID uint64) ([]*model.SunshineAward, error)
	GetStudentAwards(ctx context.Context, studentID uint64) ([]*model.SunshineAward, error)
	GetChildSunshine(ctx context.Context, childID, parentID uint64) (map[string]interface{}, error)
}

// sunshineServiceImpl 阳光章服务实现
type sunshineServiceImpl struct {
	sunshineRepo    repository.SunshineRepository
	classRepo       repository.ClassRepository
	childRepo       repository.ChildRepository
	parentChildRepo repository.ParentChildRepository
}

// NewSunshineService 创建阳光章服务实例
func NewSunshineService(
	sunshineRepo repository.SunshineRepository,
	classRepo repository.ClassRepository,
	childRepo repository.ChildRepository,
	parentChildRepo repository.ParentChildRepository,
) SunshineService {
	return &sunshineServiceImpl{
		sunshineRepo:    sunshineRepo,
		classRepo:       classRepo,
		childRepo:       childRepo,
		parentChildRepo: parentChildRepo,
	}
}

// ─── 七色配置管理 ──────────────────────────────────────────

func (s *sunshineServiceImpl) GetSunshineColors(ctx context.Context, schoolID uint64) ([]*model.SunshineColor, error) {
	// 获取学校七色-科目配置
	return s.sunshineRepo.FindColorsBySchoolID(ctx, schoolID)
}

func (s *sunshineServiceImpl) CreateSunshineColor(ctx context.Context, color *model.SunshineColor) error {
	// 配置七色-科目映射
	return s.sunshineRepo.CreateColor(ctx, color)
}

func (s *sunshineServiceImpl) UpdateSunshineColor(ctx context.Context, color *model.SunshineColor) error {
	// 更新颜色配置
	return s.sunshineRepo.UpdateColor(ctx, color)
}

// ─── 盖章操作 ──────────────────────────────────────────────

func (s *sunshineServiceImpl) StampSunshine(ctx context.Context, teacherID, studentID, colorID uint64, subject string) error {
	// 检查老师是否存在
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return err
	}
	if child == nil {
		return errors.New("学生不存在")
	}

	// 检查颜色配置是否存在
	// 这里可以添加颜色配置的检查逻辑

	// 创建盖章记录
	stamp := &model.SunshineStamp{
		StudentID: studentID,
		TeacherID: teacherID,
		ColorID:   colorID,
		Subject:   subject,
		StampDate: time.Now(),
	}

	// 保存盖章记录
	return s.sunshineRepo.CreateStamp(ctx, stamp)
}

func (s *sunshineServiceImpl) GetClassStamps(ctx context.Context, teacherID, classID uint64, params map[string]interface{}) ([]*model.SunshineStamp, int64, error) {
	// 检查班级是否存在
	// 由于缺少ClassRepository的具体方法，暂时跳过

	// 解析分页参数
	limit := config.Get().Behavior.DefaultLimit
	offset := config.Get().Behavior.DefaultOffset
	if l, ok := params["limit"].(int); ok && l > 0 {
		limit = l
	}
	if page, ok := params["page"].(int); ok && page > 0 {
		offset = (page - 1) * limit
	}

	// 构建查询条件
	// 由于缺少班级与学生的关联关系，暂时跳过班级过滤

	// 获取班级盖章记录
	return s.sunshineRepo.FindStampsByClassID(ctx, classID, params, limit, offset)
}

func (s *sunshineServiceImpl) GetStudentStamps(ctx context.Context, studentID uint64) (map[string]interface{}, error) {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("学生不存在")
	}

	// 获取学生的盖章记录
	stamps, err := s.sunshineRepo.FindStampsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// 统计各颜色的盖章数量
	stampCount := make(map[string]int)
	for _, stamp := range stamps {
		// 根据colorID获取颜色名称
		colorName := "未知"
		if name, ok := config.Get().Sunshine.ColorMap[stamp.ColorID]; ok {
			colorName = name
		}
		stampCount[colorName]++
	}

	return map[string]interface{}{"stamps": stampCount}, nil
}

// ─── 之星评选 ──────────────────────────────────────────────

func (s *sunshineServiceImpl) EvaluateSunshineAwards(ctx context.Context, teacherID, classID uint64, period string) error {
	// 检查班级是否存在
	// 由于缺少ClassRepository的具体方法，暂时跳过

	// 触发评选逻辑
	// 1. 获取班级学生列表
	// 2. 统计每个学生的盖章数量
	// 3. 为每个颜色选出盖章最多的学生
	// 4. 创建评选记录

	// 由于缺少班级学生列表的获取方法，暂时返回nil
	return nil
}

func (s *sunshineServiceImpl) GetSunshineAwards(ctx context.Context, teacherID, classID uint64) ([]*model.SunshineAward, error) {
	// 检查班级是否存在
	// 由于缺少ClassRepository的具体方法，暂时跳过

	// 获取评选结果
	return s.sunshineRepo.FindAwardsByClassID(ctx, classID)
}

func (s *sunshineServiceImpl) GetStudentAwards(ctx context.Context, studentID uint64) ([]*model.SunshineAward, error) {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("学生不存在")
	}

	// 获取学生的之星称号
	return s.sunshineRepo.FindAwardsByStudentID(ctx, studentID)
}

func (s *sunshineServiceImpl) GetChildSunshine(ctx context.Context, childID, parentID uint64) (map[string]interface{}, error) {
	// 检查家长是否有权限查看该孩子
	bindings, err := s.parentChildRepo.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// 检查是否存在绑定关系
	found := false
	for _, binding := range bindings {
		if binding.ChildID == childID {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("您没有权限查看该孩子的信息")
	}

	// 检查孩子是否存在
	child, err := s.childRepo.FindByID(ctx, childID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("孩子不存在")
	}

	// 获取孩子的盖章记录
	stamps, err := s.sunshineRepo.FindStampsByStudentID(ctx, childID)
	if err != nil {
		return nil, err
	}

	// 统计各颜色的盖章数量
	stampCount := make(map[string]int)
	for _, stamp := range stamps {
		// 根据colorID获取颜色名称
		colorName := "未知"
		if name, ok := config.Get().Sunshine.ColorMap[stamp.ColorID]; ok {
			colorName = name
		}
		stampCount[colorName]++
	}

	// 获取孩子的之星称号
	awards, err := s.sunshineRepo.FindAwardsByStudentID(ctx, childID)
	if err != nil {
		return nil, err
	}

	// 转换奖项格式
	awardList := make([]map[string]interface{}, len(awards))
	for i, award := range awards {
		awardList[i] = map[string]interface{}{
			"id":           award.ID,
			"award_name":   award.AwardName,
			"period":       award.Period,
			"period_year":  award.PeriodYear,
			"period_month": award.PeriodMonth,
		}
	}

	return map[string]interface{}{
		"stamps": stampCount,
		"awards": awardList,
	}, nil
}

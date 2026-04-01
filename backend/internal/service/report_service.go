// growth-partner/backend/internal/service/report_service.go
// 报告服务：生成周报、月度卡片及年度成长画卷（PDF/图片导出）

package service

import (
	"context"
	"fmt"
	"growth-partner/internal/repository"
)

type ReportService interface {
	GenerateWeeklyReport(ctx context.Context, childID uint64) (string, error)
}

type reportServiceImpl struct {
	behaviorRepo repository.BehaviorRepository
	childRepo    repository.ChildRepository
}

func NewReportService(b repository.BehaviorRepository, c repository.ChildRepository) ReportService {
	return &reportServiceImpl{behaviorRepo: b, childRepo: c}
}

func (s *reportServiceImpl) GenerateWeeklyReport(ctx context.Context, childID uint64) (string, error) {
	// 这里未来会调用 chromedp 渲染 HTML 模板并输出 PDF
	// 目前仅做逻辑占位
	return fmt.Sprintf("/reports/weekly_%d.pdf", childID), nil
}

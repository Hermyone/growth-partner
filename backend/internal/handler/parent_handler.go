// growth-partner/backend/internal/handler/parent_handler.go
// 家长端模块控制器

package handler

import (
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// ParentHandler 家长端控制器
type ParentHandler struct {
	parentSvc service.ParentService
}

// NewParentHandler 创建家长端控制器实例
func NewParentHandler(parentSvc service.ParentService) *ParentHandler {
	return &ParentHandler{
		parentSvc: parentSvc,
	}
}

// ─── 家长端核心功能 ──────────────────────────────────────────

// GET /api/v1/parent/children
func (h *ParentHandler) GetMyChildren(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取绑定的孩子列表
	children, err := h.parentSvc.GetMyChildren(c.Request.Context(), parentID.(uint64))
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, children)
}

// GET /api/v1/parent/children/:childId/partner
func (h *ParentHandler) GetChildPartner(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 调用服务获取孩子的伙伴
	partner, err := h.parentSvc.GetChildPartner(c.Request.Context(), childID, parentID.(uint64))
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, partner)
}

// GET /api/v1/parent/children/:childId/partners
func (h *ParentHandler) GetChildPartners(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 调用服务获取孩子的历史伙伴列表
	partners, err := h.parentSvc.GetChildPartners(c.Request.Context(), childID, parentID.(uint64))
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, partners)
}

// GET /api/v1/parent/children/:childId/behaviors
func (h *ParentHandler) GetChildBehaviors(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 构建查询参数
	params := map[string]interface{}{}
	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		params["page"] = page
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		params["limit"] = limit
	}
	if startDate := c.Query("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		params["end_date"] = endDate
	}
	if dimension := c.Query("dimension"); dimension != "" {
		params["dimension"] = dimension
	}

	// 调用服务获取孩子的行为记录
	behaviors, total, err := h.parentSvc.GetChildBehaviors(c.Request.Context(), childID, parentID.(uint64), params)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{
		"data":  behaviors,
		"total": total,
	})
}

// GET /api/v1/parent/children/:childId/broadcasts
func (h *ParentHandler) GetChildBroadcasts(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 调用服务获取孩子收到的广播
	broadcasts, err := h.parentSvc.GetChildBroadcasts(c.Request.Context(), childID, parentID.(uint64))
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, broadcasts)
}

// GET /api/v1/parent/children/:childId/milestones
func (h *ParentHandler) GetChildMilestones(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 调用服务获取孩子的里程碑
	milestones, err := h.parentSvc.GetChildMilestones(c.Request.Context(), childID, parentID.(uint64))
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, milestones)
}

// GET /api/v1/parent/children/:childId/monthly-card
func (h *ParentHandler) GetChildMonthlyCard(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 获取月份参数
	month := c.Query("month")
	if month == "" {
		middleware.ResponseValidationError(c, "无效的月份参数")
		return
	}

	// 调用服务获取孩子的月度成长卡
	card, err := h.parentSvc.GetChildMonthlyCard(c.Request.Context(), childID, parentID.(uint64), month)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, card)
}

// GET /api/v1/parent/children/:childId/annual-report
func (h *ParentHandler) GetChildAnnualReport(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 获取年份参数
	year := c.Query("year")
	if year == "" {
		middleware.ResponseValidationError(c, "无效的年份参数")
		return
	}

	// 调用服务获取孩子的年度成长画卷
	report, err := h.parentSvc.GetChildAnnualReport(c.Request.Context(), childID, parentID.(uint64), year)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, report)
}

// GET /api/v1/parent/children/:childId/battles
func (h *ParentHandler) GetChildBattles(c *gin.Context) {
	// 从上下文获取家长ID
	parentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取孩子ID
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的孩子ID")
		return
	}

	// 调用服务获取孩子的对战参与记录
	battles, err := h.parentSvc.GetChildBattles(c.Request.Context(), childID, parentID.(uint64))
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, battles)
}

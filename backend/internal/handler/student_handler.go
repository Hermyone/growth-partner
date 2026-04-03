// growth-partner/backend/internal/handler/student_handler.go
// 学生端模块控制器

package handler

import (
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// StudentHandler 学生端控制器
type StudentHandler struct {
	studentSvc service.StudentService
}

// NewStudentHandler 创建学生端控制器实例
func NewStudentHandler(studentSvc service.StudentService) *StudentHandler {
	return &StudentHandler{
		studentSvc: studentSvc,
	}
}

// ─── 4.1 伙伴系统 ──────────────────────────────────────────

// GET /api/v1/student/partner
func (h *StudentHandler) GetCurrentPartner(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取当前伙伴
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	partner, err := h.studentSvc.GetCurrentPartner(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, partner)
}

// GET /api/v1/student/partners
func (h *StudentHandler) GetAllPartners(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取所有伙伴
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	partners, err := h.studentSvc.GetAllPartners(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, partners)
}

// POST /api/v1/student/partner
func (h *StudentHandler) SelectNewPartner(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type SelectNewPartnerReq struct {
		TemplateID uint64 `json:"template_id" binding:"required"`
	}

	var req SelectNewPartnerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务选择新伙伴
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.studentSvc.SelectNewPartner(c.Request.Context(), studentIDUint, req.TemplateID)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "选择伙伴成功"})
}

// PATCH /api/v1/student/partner/nickname
func (h *StudentHandler) UpdatePartnerNickname(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type UpdatePartnerNicknameReq struct {
		Nickname string `json:"nickname" binding:"required"`
	}

	var req UpdatePartnerNicknameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务更新伙伴昵称
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.studentSvc.UpdatePartnerNickname(c.Request.Context(), studentIDUint, req.Nickname)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "更新昵称成功"})
}

// GET /api/v1/student/partner/growth-history
func (h *StudentHandler) GetPartnerGrowthHistory(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
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

	// 调用服务获取成长值流水
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	records, total, err := h.studentSvc.GetPartnerGrowthHistory(c.Request.Context(), studentIDUint, params)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{
		"data":  records,
		"total": total,
	})
}

// GET /api/v1/student/partner/templates
func (h *StudentHandler) GetPartnerTemplates(c *gin.Context) {
	// 调用服务获取伙伴模板列表
	templates, err := h.studentSvc.GetPartnerTemplates(c.Request.Context())
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, templates)
}

// ─── 4.2 行为记录查看 ───────────────────────────────────────

// GET /api/v1/student/behaviors
func (h *StudentHandler) GetMyBehaviors(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
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

	// 调用服务获取行为记录
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	behaviors, total, err := h.studentSvc.GetMyBehaviors(c.Request.Context(), studentIDUint, params)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{
		"data":  behaviors,
		"total": total,
	})
}

// GET /api/v1/student/behaviors/stats
func (h *StudentHandler) GetBehaviorStats(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取行为统计数据
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	stats, err := h.studentSvc.GetBehaviorStats(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, stats)
}

// ─── 4.3 广播中心 ──────────────────────────────────────────

// GET /api/v1/student/broadcasts
func (h *StudentHandler) GetMyBroadcasts(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取广播列表
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	broadcasts, err := h.studentSvc.GetMyBroadcasts(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, broadcasts)
}

// PATCH /api/v1/student/broadcasts/:id/read
func (h *StudentHandler) MarkBroadcastAsRead(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取广播ID
	broadcastID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的广播ID")
		return
	}

	// 调用服务标记广播为已读
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.studentSvc.MarkBroadcastAsRead(c.Request.Context(), broadcastID, studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "标记已读成功"})
}

// POST /api/v1/student/broadcasts/read-all
func (h *StudentHandler) MarkAllBroadcastsAsRead(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务一键标记所有广播为已读
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.studentSvc.MarkAllBroadcastsAsRead(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "标记所有已读成功"})
}

// ─── 4.4 成长年历 ──────────────────────────────────────────

// GET /api/v1/student/growth-calendar/months
func (h *StudentHandler) GetGrowthCalendarMonths(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取月度成长卡列表
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	months, err := h.studentSvc.GetGrowthCalendarMonths(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, months)
}

// GET /api/v1/student/growth-calendar/months/:month
func (h *StudentHandler) GetGrowthCalendarMonth(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取月份参数
	month := c.Param("month")
	if month == "" {
		middleware.ResponseValidationError(c, "无效的月份参数")
		return
	}

	// 调用服务获取指定月份的成长卡详情
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	card, err := h.studentSvc.GetGrowthCalendarMonth(c.Request.Context(), studentIDUint, month)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, card)
}

// GET /api/v1/student/growth-calendar/annual/:year
func (h *StudentHandler) GetGrowthCalendarAnnual(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取年份参数
	year := c.Param("year")
	if year == "" {
		middleware.ResponseValidationError(c, "无效的年份参数")
		return
	}

	// 调用服务获取年度成长画卷数据
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	report, err := h.studentSvc.GetGrowthCalendarAnnual(c.Request.Context(), studentIDUint, year)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, report)
}

// GET /api/v1/student/milestones
func (h *StudentHandler) GetMilestones(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取里程碑列表
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	milestones, err := h.studentSvc.GetMilestones(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, milestones)
}

// ─── 4.5 盲盒查看 ──────────────────────────────────────────

// GET /api/v1/student/blindbox/my-draws
func (h *StudentHandler) GetMyBlindboxDraws(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取盲盒抽取记录
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	draws, err := h.studentSvc.GetMyBlindboxDraws(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, draws)
}

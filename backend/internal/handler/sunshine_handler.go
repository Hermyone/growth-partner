// growth-partner/backend/internal/handler/sunshine_handler.go
// 阳光章模块控制器

package handler

import (
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/model"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// SunshineHandler 阳光章控制器
type SunshineHandler struct {
	sunshineSvc service.SunshineService
}

// NewSunshineHandler 创建阳光章控制器实例
func NewSunshineHandler(sunshineSvc service.SunshineService) *SunshineHandler {
	return &SunshineHandler{
		sunshineSvc: sunshineSvc,
	}
}

// ─── 七色配置管理 ──────────────────────────────────────────

// GET /api/v1/admin/sunshine/colors
func (h *SunshineHandler) GetSunshineColors(c *gin.Context) {
	// 从上下文获取学校ID
	schoolID, err := strconv.ParseUint(c.Query("school_id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的学校ID")
		return
	}

	// 调用服务获取学校七色-科目配置
	colors, err := h.sunshineSvc.GetSunshineColors(c.Request.Context(), schoolID)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, colors)
}

// POST /api/v1/admin/sunshine/colors
func (h *SunshineHandler) CreateSunshineColor(c *gin.Context) {
	// 定义请求结构
	var color model.SunshineColor
	if err := c.ShouldBindJSON(&color); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务配置七色-科目映射
	err := h.sunshineSvc.CreateSunshineColor(c.Request.Context(), &color)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "创建阳光章颜色配置成功"})
}

// PUT /api/v1/admin/sunshine/colors/:id
func (h *SunshineHandler) UpdateSunshineColor(c *gin.Context) {
	// 获取颜色ID
	colorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的颜色ID")
		return
	}

	// 定义请求结构
	var color model.SunshineColor
	if err := c.ShouldBindJSON(&color); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 设置颜色ID
	color.ID = colorID

	// 调用服务更新颜色配置
	err = h.sunshineSvc.UpdateSunshineColor(c.Request.Context(), &color)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "更新阳光章颜色配置成功"})
}

// ─── 盖章操作 ──────────────────────────────────────────────

// POST /api/v1/teacher/sunshine/stamp
func (h *SunshineHandler) StampSunshine(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type StampSunshineReq struct {
		StudentID uint64 `json:"student_id" binding:"required"`
		ColorID   uint64 `json:"color_id" binding:"required"`
		Subject   string `json:"subject" binding:"required"`
	}

	var req StampSunshineReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务为学生盖章
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.sunshineSvc.StampSunshine(c.Request.Context(), teacherIDUint, req.StudentID, req.ColorID, req.Subject)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "盖章成功"})
}

// GET /api/v1/teacher/sunshine/stamps
func (h *SunshineHandler) GetClassStamps(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取班级ID
	classID, err := strconv.ParseUint(c.Query("class_id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 构建查询参数
	params := map[string]interface{}{}
	if month := c.Query("month"); month != "" {
		params["month"] = month
	}
	if color := c.Query("color"); color != "" {
		params["color"] = color
	}

	// 调用服务查看班级盖章记录
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	stamps, total, err := h.sunshineSvc.GetClassStamps(c.Request.Context(), teacherIDUint, classID, params)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{
		"data":  stamps,
		"total": total,
	})
}

// GET /api/v1/student/sunshine/my-stamps
func (h *SunshineHandler) GetStudentStamps(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务查看学生盖章概况
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	stamps, err := h.sunshineSvc.GetStudentStamps(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, stamps)
}

// ─── 之星评选 ──────────────────────────────────────────────

// POST /api/v1/teacher/sunshine/awards/evaluate
func (h *SunshineHandler) EvaluateSunshineAwards(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type EvaluateSunshineAwardsReq struct {
		ClassID uint64 `json:"class_id" binding:"required"`
		Period  string `json:"period" binding:"required"` // 月度/季度/年度
	}

	var req EvaluateSunshineAwardsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务触发之星评选
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.sunshineSvc.EvaluateSunshineAwards(c.Request.Context(), teacherIDUint, req.ClassID, req.Period)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "评选成功"})
}

// GET /api/v1/teacher/sunshine/awards
func (h *SunshineHandler) GetSunshineAwards(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取班级ID
	classID, err := strconv.ParseUint(c.Query("class_id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 调用服务查看评选结果列表
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	awards, err := h.sunshineSvc.GetSunshineAwards(c.Request.Context(), teacherIDUint, classID)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, awards)
}

// GET /api/v1/student/sunshine/my-awards
func (h *SunshineHandler) GetStudentAwards(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务查看学生之星称号
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	awards, err := h.sunshineSvc.GetStudentAwards(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, awards)
}

// GET /api/v1/parent/children/:childId/sunshine
func (h *SunshineHandler) GetChildSunshine(c *gin.Context) {
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

	// 调用服务查看孩子阳光章情况
	parentIDUint, ok := parentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	sunshine, err := h.sunshineSvc.GetChildSunshine(c.Request.Context(), childID, parentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, sunshine)
}

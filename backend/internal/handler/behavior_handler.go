// growth-partner/backend/internal/handler/behavior_handler.go
// 行为记录接口层：解析 HTTP 请求，调用服务层，返回统一格式响应

package handler

import (
	"net/http"
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/model"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// BehaviorHandler 行为记录控制器
type BehaviorHandler struct {
	behaviorSvc service.BehaviorService
}

// NewBehaviorHandler 创建行为记录控制器
func NewBehaviorHandler(behaviorSvc service.BehaviorService) *BehaviorHandler {
	return &BehaviorHandler{
		behaviorSvc: behaviorSvc,
	}
}

// RecordBehaviorReq 记录行为的请求参数定义
type RecordBehaviorReq struct {
	ChildID     uint64                  `json:"child_id" binding:"required"`
	Dimension   model.BehaviorDimension `json:"dimension" binding:"required"`
	Description string                  `json:"description" binding:"required,max=256"`
	GrowthValue int                     `json:"growth_value" binding:"required,min=1,max=10"`
}

// RecordBehavior 录入学生正向行为 (老师/专业教练调用)
// POST /api/v1/teacher/behaviors
func (h *BehaviorHandler) RecordBehavior(c *gin.Context) {
	var req RecordBehaviorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 从 JWT Token 中提取当前操作人的信息
	recorderID := middleware.GetUserID(c)
	recorderRole := middleware.GetRole(c)
	classID := middleware.GetClassID(c) // 老师所在的班级

	// 构造传给 Service 的请求
	svcReq := service.RecordBehaviorRequest{
		ChildID:      req.ChildID,
		ClassID:      classID,
		RecorderID:   recorderID,
		RecorderRole: model.RecorderRole(recorderRole),
		Dimension:    req.Dimension,
		Description:  req.Description,
		GrowthValue:  req.GrowthValue,
	}

	// 调用业务逻辑
	record, growthResult, err := h.behaviorSvc.RecordBehavior(c.Request.Context(), svcReq)
	if err != nil {
		middleware.ResponseError(c, http.StatusBadRequest, "RECORD_BEHAVIOR_FAILED", err.Error())
		return
	}

	// 返回成功结果
	response := gin.H{
		"behavior_record": record,
	}
	
	if growthResult != nil {
		response["growth_result"] = gin.H{
			"is_evolved":      growthResult.IsEvolved,
			"current_points":  growthResult.Partner.GrowthPoints,
			"partner_message": growthResult.PartnerMessage,
			"evolution_msg":   growthResult.EvolutionMsg,
		}
	}
	
	middleware.ResponseOK(c, response)
}

// GetClassBehaviors 获取班级行为动态 (老师/园长调用)
// GET /api/v1/teacher/behaviors?page=1&size=20
func (h *BehaviorHandler) GetClassBehaviors(c *gin.Context) {
	classID := middleware.GetClassID(c)
	page, size := getPaginationParams(c)
	offset := (page - 1) * size

	records, total, err := h.behaviorSvc.GetClassBehaviors(c.Request.Context(), classID, size, offset)
	if err != nil {
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  records,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetChildBehaviors 获取单个学生的行为动态 (家长/学生端调用)
// GET /api/v1/parent/children/:childId/behaviors?page=1&size=20
func (h *BehaviorHandler) GetChildBehaviors(c *gin.Context) {
	childIDStr := c.Param("childId")
	childID, err := strconv.ParseUint(childIDStr, 10, 64)
	if err != nil {
		middleware.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "无效的学生ID")
		return
	}

	// TODO: 这里应当有一步鉴权，确保当前家长只能查看自己绑定的 childID

	page, size := getPaginationParams(c)
	offset := (page - 1) * size

	records, total, err := h.behaviorSvc.GetChildBehaviors(c.Request.Context(), childID, size, offset)
	if err != nil {
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  records,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

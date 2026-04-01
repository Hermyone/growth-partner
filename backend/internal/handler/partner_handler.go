// growth-partner/backend/internal/handler/partner_handler.go

package handler

import (
	"net/http"
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// PartnerHandler 伙伴控制器
type PartnerHandler struct {
	partnerSvc service.PartnerService
}

// NewPartnerHandler 创建伙伴控制器
func NewPartnerHandler(partnerSvc service.PartnerService) *PartnerHandler {
	return &PartnerHandler{
		partnerSvc: partnerSvc,
	}
}

// ─── 公开接口：伙伴模板 ───────────────────────────────────────

// ListTemplates 获取所有可用的伙伴模板
// GET /api/v1/partner-templates
func (h *PartnerHandler) ListTemplates(c *gin.Context) {
	templates, err := h.partnerSvc.GetAllTemplates(c.Request.Context())
	if err != nil {
		middleware.ResponseInternalError(c, true, err)
		return
	}
	middleware.ResponseOK(c, templates)
}

// GetTemplate 获取单个模板详情
// GET /api/v1/partner-templates/:id
func (h *PartnerHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "无效的模板ID")
		return
	}

	template, err := h.partnerSvc.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			middleware.ResponseError(c, http.StatusNotFound, "NOT_FOUND", "模板不存在")
			return
		}
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOK(c, template)
}

// ─── 学生端接口：专属伙伴 ─────────────────────────────────────

// CreatePartnerReq 创建伙伴请求
type CreatePartnerReq struct {
	TemplateID uint64 `json:"template_id" binding:"required"`
	Nickname   string `json:"nickname" binding:"required,max=32"`
}

// CreatePartner 领养专属伙伴（每个学生仅限一次）
// POST /api/v1/partner
func (h *PartnerHandler) CreatePartner(c *gin.Context) {
	var req CreatePartnerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	childID := middleware.GetChildID(c)
	if childID == 0 {
		middleware.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "仅限学生账号操作")
		return
	}

	partner, err := h.partnerSvc.CreatePartner(c.Request.Context(), childID, req.TemplateID, req.Nickname)
	if err != nil {
		if err == service.ErrPartnerAlreadyExist {
			middleware.ResponseError(c, http.StatusConflict, "ALREADY_EXISTS", "你已经拥有一个伙伴啦，不能重复领养哦")
			return
		}
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOKWithMessage(c, "领养伙伴成功！", partner)
}

// GetMyPartner 获取我的伙伴状态
// GET /api/v1/partner
func (h *PartnerHandler) GetMyPartner(c *gin.Context) {
	childID := middleware.GetChildID(c)
	partner, err := h.partnerSvc.GetPartnerByChildID(c.Request.Context(), childID)
	if err != nil {
		if err == service.ErrPartnerNotFound {
			middleware.ResponseError(c, http.StatusNotFound, "NOT_FOUND", "你还没有领养伙伴哦")
			return
		}
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOK(c, partner)
}

// UpdateNicknameReq 修改昵称请求
type UpdateNicknameReq struct {
	Nickname string `json:"nickname" binding:"required,max=32"`
}

// UpdateNickname 修改伙伴昵称
// PATCH /api/v1/partner/nickname
func (h *PartnerHandler) UpdateNickname(c *gin.Context) {
	var req UpdateNicknameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	childID := middleware.GetChildID(c)
	err := h.partnerSvc.UpdateNickname(c.Request.Context(), childID, req.Nickname)
	if err != nil {
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOKWithMessage(c, "昵称修改成功", nil)
}

// GetGrowthHistory 获取我的伙伴成长历史
// GET /api/v1/partner/growth-history?page=1&size=20
func (h *PartnerHandler) GetGrowthHistory(c *gin.Context) {
	childID := middleware.GetChildID(c)
	page, size := getPaginationParams(c)
	offset := (page - 1) * size

	records, total, err := h.partnerSvc.GetGrowthHistory(c.Request.Context(), childID, size, offset)
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

// ─── 教师/家长端接口：查看学生的伙伴 ─────────────────────────────

// GetChildPartner 获取指定学生的伙伴（家长/老师使用）
// GET /api/v1/parent/children/:childId/partner
func (h *PartnerHandler) GetChildPartner(c *gin.Context) {
	childIDStr := c.Param("childId")
	childID, err := strconv.ParseUint(childIDStr, 10, 64)
	if err != nil {
		middleware.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "无效的学生ID")
		return
	}

	// 这里可以加入老师/家长对该 childID 的权限校验

	partner, err := h.partnerSvc.GetPartnerByChildID(c.Request.Context(), childID)
	if err != nil {
		if err == service.ErrPartnerNotFound {
			middleware.ResponseError(c, http.StatusNotFound, "NOT_FOUND", "该学生尚未领养伙伴")
			return
		}
		middleware.ResponseInternalError(c, true, err)
		return
	}

	middleware.ResponseOK(c, partner)
}

// GetClassOverview 获取全班伙伴概览（班主任使用，只显示进度不显示排名）
// GET /api/v1/teacher/classes/:classId/overview
func (h *PartnerHandler) GetClassOverview(c *gin.Context) {
	// 暂留接口，调用相应的 service 方法聚合班级数据
	middleware.ResponseOK(c, gin.H{"message": "获取班级概览成功（开发中）"})
}

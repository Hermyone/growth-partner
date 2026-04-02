// growth-partner/backend/internal/handler/partner_template_handler.go
// 伙伴模板模块控制器

package handler

import (
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/model"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// PartnerTemplateHandler 伙伴模板控制器
type PartnerTemplateHandler struct {
	templateSvc service.PartnerTemplateService
}

// NewPartnerTemplateHandler 创建伙伴模板控制器实例
func NewPartnerTemplateHandler(templateSvc service.PartnerTemplateService) *PartnerTemplateHandler {
	return &PartnerTemplateHandler{
		templateSvc: templateSvc,
	}
}

// ─── 公开接口 ───────────────────────────────────────────────

// GET /api/v1/partner-templates
func (h *PartnerTemplateHandler) GetAllPartnerTemplates(c *gin.Context) {
	// 调用服务获取全部激活的伙伴模板列表
	templates, err := h.templateSvc.GetAllActivePartnerTemplates(c.Request.Context())
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, templates)
}

// GET /api/v1/partner-templates/:id
func (h *PartnerTemplateHandler) GetPartnerTemplateByID(c *gin.Context) {
	// 获取模板ID
	templateID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的模板ID")
		return
	}

	// 调用服务获取单个伙伴模板详情
	template, err := h.templateSvc.GetPartnerTemplateByID(c.Request.Context(), templateID)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, template)
}

// GET /api/v1/health
func (h *PartnerTemplateHandler) GetHealthStatus(c *gin.Context) {
	// 调用服务获取健康状态
	status := h.templateSvc.GetHealthStatus(c.Request.Context())
	middleware.ResponseOK(c, status)
}

// GET /api/v1/config/client
func (h *PartnerTemplateHandler) GetClientConfig(c *gin.Context) {
	// 调用服务获取前端全局配置
	config := h.templateSvc.GetClientConfig(c.Request.Context())
	middleware.ResponseOK(c, config)
}

// ─── 管理员接口 ─────────────────────────────────────────────

// POST /api/v1/admin/partner-templates
func (h *PartnerTemplateHandler) CreatePartnerTemplate(c *gin.Context) {
	// 定义请求结构
	var template model.PartnerTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务新增伙伴模板
	err := h.templateSvc.CreatePartnerTemplate(c.Request.Context(), &template)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "创建伙伴模板成功"})
}

// PUT /api/v1/admin/partner-templates/:id
func (h *PartnerTemplateHandler) UpdatePartnerTemplate(c *gin.Context) {
	// 获取模板ID
	templateID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的模板ID")
		return
	}

	// 定义请求结构
	var template model.PartnerTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 设置模板ID
	template.ID = templateID

	// 调用服务更新模板信息
	err = h.templateSvc.UpdatePartnerTemplate(c.Request.Context(), &template)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "更新伙伴模板成功"})
}

// POST /api/v1/admin/partner-templates/seed
func (h *PartnerTemplateHandler) SeedPartnerTemplates(c *gin.Context) {
	// 调用服务一键初始化30个预设模板
	err := h.templateSvc.SeedPartnerTemplates(c.Request.Context())
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, gin.H{"message": "初始化伙伴模板成功"})
}

// growth-partner/backend/internal/router/partner_template_router.go
// 伙伴模板路由

package router

import (
	"growth-partner/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterPartnerTemplateRoutes 注册伙伴模板路由
func RegisterPartnerTemplateRoutes(router *gin.RouterGroup, templateHandler *handler.PartnerTemplateHandler) {
	// 公开接口
	router.GET("/partner-templates", templateHandler.GetAllPartnerTemplates)
	router.GET("/partner-templates/:id", templateHandler.GetPartnerTemplateByID)
	router.GET("/health", templateHandler.GetHealthStatus)
	router.GET("/config/client", templateHandler.GetClientConfig)

	// 注意：管理员接口已在 admin_router.go 中注册，避免重复注册
}

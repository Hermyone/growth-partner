// growth-partner/backend/internal/router/partner_template_router.go
// 伙伴模板路由

package router

import (
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPartnerTemplateRoutes 注册伙伴模板路由
func RegisterPartnerTemplateRoutes(router *gin.RouterGroup, templateHandler *handler.PartnerTemplateHandler) {
	// 公开接口
	router.GET("/partner-templates", templateHandler.GetAllPartnerTemplates)
	router.GET("/partner-templates/:id", templateHandler.GetPartnerTemplateByID)
	router.GET("/health", templateHandler.GetHealthStatus)
	router.GET("/config/client", templateHandler.GetClientConfig)

	// 管理员接口
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.RequireRoles("admin"))
	{
		adminRoutes.POST("/partner-templates", templateHandler.CreatePartnerTemplate)
		adminRoutes.PUT("/partner-templates/:id", templateHandler.UpdatePartnerTemplate)
		adminRoutes.POST("/partner-templates/seed", templateHandler.SeedPartnerTemplates)
	}
}

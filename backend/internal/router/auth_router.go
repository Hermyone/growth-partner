// growth-partner/backend/internal/router/auth_router.go
// 认证模块路由配置

package router

import (
	"growth-partner/config"
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"
	"growth-partner/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 配置认证相关路由
func SetupAuthRoutes(
	v1 *gin.RouterGroup,
	cfg *config.Config,
	jwtManager *jwt.Manager,
	authHandler *handler.AuthHandler,
) {
	// ─── 公开接口（无需登录）──────────────────────────────────
	public := v1.Group("")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
	}

	// ─── 需要认证的接口（通用）───────────────────────────────
	authRequired := v1.Group("", middleware.Auth(jwtManager))
	{
		// 认证相关接口
		auth := authRequired.Group("/auth")
		{
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", authHandler.Me)
			auth.PATCH("/password", authHandler.ChangePassword)
		}
	}
}

// growth-partner/backend/internal/router/router.go
// 路由注册：清晰分层，按角色分组，统一应用中间件

package router

import (
	"net/http"

	"growth-partner/config"
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"
	"growth-partner/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// SetupRouter 创建并配置 Gin 路由
func SetupRouter(
	cfg *config.Config,
	jwtManager *jwt.Manager,
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	teacherHandler *handler.TeacherHandler,
	studentHandler *handler.StudentHandler,
	parentHandler *handler.ParentHandler,
	partnerHandler *handler.PartnerHandler,
	behaviorHandler *handler.BehaviorHandler,
	classHandler *handler.ClassHandler,
	childHandler *handler.ChildHandler,
	broadcastHandler *handler.BroadcastHandler,
	battleHandler *handler.BattleHandler,
	blindboxHandler *handler.BlindboxHandler,
	wsHandler *handler.WebSocketHandler,
	templateHandler *handler.PartnerTemplateHandler,
	sunshineHandler *handler.SunshineHandler,
) *gin.Engine {

	// 生产环境关闭调试日志
	if !cfg.App.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// ─── 全局中间件 ───────────────────────────────────────────
	r.Use(gin.Recovery())                   // panic 恢复
	r.Use(middleware.RequestLogger())       // 请求日志
	r.Use(middleware.RequestID())           // 请求追踪ID
	r.Use(middleware.CORS(cfg.App.IsDev())) // 跨域配置

	// ─── 健康检查（无需认证）─────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": config.Get().App.Name,
			"version": config.Get().App.Version,
		})
	})

	// ─── API v1 路由组 ────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// ─── 配置各模块路由 ───────────────────────────────────────
	SetupAuthRoutes(v1, cfg, jwtManager, authHandler)
	SetupAdminRoutes(v1, cfg, jwtManager, adminHandler)
	SetupTeacherRoutes(v1, cfg, jwtManager, teacherHandler, behaviorHandler, broadcastHandler, blindboxHandler)
	RegisterStudentRoutes(v1, studentHandler)
	RegisterParentRoutes(v1, parentHandler)
	RegisterBattleRoutes(v1, battleHandler)
	RegisterWebSocketRoutes(v1, wsHandler, jwtManager)
	RegisterPartnerTemplateRoutes(v1, templateHandler)
	RegisterSunshineRoutes(v1, sunshineHandler)

	// ─── 公开接口（无需登录）──────────────────────────────────
	public := v1.Group("")
	{
		// 伙伴模板（公开，供选伙伴时展示）
		public.GET("/partner-templates", partnerHandler.ListTemplates)
		public.GET("/partner-templates/:id", partnerHandler.GetTemplate)
	}

	return r
}

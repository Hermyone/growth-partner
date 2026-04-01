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
	partnerHandler *handler.PartnerHandler,
	behaviorHandler *handler.BehaviorHandler,
	classHandler *handler.ClassHandler,
	childHandler *handler.ChildHandler,
	broadcastHandler *handler.BroadcastHandler,
	battleHandler *handler.BattleHandler,
	blindboxHandler *handler.BlindboxHandler,
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
			"service": "成长伙伴",
			"version": "1.0.0",
		})
	})

	// ─── API v1 路由组 ────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// ─── 公开接口（无需登录）──────────────────────────────────
	public := v1.Group("")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
		// 伙伴模板（公开，供选伙伴时展示）
		public.GET("/partner-templates", partnerHandler.ListTemplates)
		public.GET("/partner-templates/:id", partnerHandler.GetTemplate)
	}

	// ─── 需要认证的接口（通用）───────────────────────────────
	authRequired := v1.Group("", middleware.Auth(jwtManager))
	{
		// 伙伴接口（学生可用）
		partnerGroup := authRequired.Group("/partner")
		{
			partnerGroup.GET("", middleware.RequireStudent(), partnerHandler.GetMyPartner)
			partnerGroup.POST("", middleware.RequireStudent(), partnerHandler.CreatePartner)
			partnerGroup.PATCH("/nickname", middleware.RequireStudent(), partnerHandler.UpdateNickname)
			partnerGroup.GET("/growth-history", middleware.RequireStudent(), partnerHandler.GetGrowthHistory)
		}

		// 对战接口（学生可用）
		battleGroup := authRequired.Group("/battle")
		{
			battleGroup.POST("/room", middleware.RequireStudent(), battleHandler.CreateRoom)
			battleGroup.POST("/room/:roomId/join", middleware.RequireStudent(), battleHandler.JoinRoom)
			battleGroup.GET("/ws", middleware.RequireStudent(), battleHandler.WebSocketUpgrade)
			battleGroup.GET("/history", middleware.RequireStudent(), battleHandler.GetMyHistory)
		}

		// 广播接口
		broadcastGroup := authRequired.Group("/broadcast")
		{
			broadcastGroup.GET("/ws", broadcastHandler.WebSocketUpgrade) // WebSocket 连接
			broadcastGroup.GET("/list", broadcastHandler.GetBroadcastList)
			broadcastGroup.PATCH("/:id/read", broadcastHandler.MarkAsRead)
			// 园长发送广播（仅教师）
			broadcastGroup.POST("", middleware.RequireTeacher(), broadcastHandler.Send)
		}

		// 家长端接口
		parentGroup := authRequired.Group("/parent", middleware.RequireParent())
		{
			parentGroup.GET("/children", childHandler.GetMyChildren)
			parentGroup.GET("/children/:childId/partner", partnerHandler.GetChildPartner)
			parentGroup.GET("/children/:childId/behaviors", behaviorHandler.GetChildBehaviors)
		}
	}

	// ─── 教师/园长端接口 ─────────────────────────────────────
	teacherAPI := v1.Group("/teacher",
		middleware.Auth(jwtManager),
		middleware.RequireTeacher(),
	)
	{
		// 班级管理
		teacherAPI.GET("/classes", classHandler.GetMyClasses)
		teacherAPI.POST("/classes", classHandler.CreateClass)
		teacherAPI.GET("/classes/:classId/students", childHandler.GetClassStudents)
		teacherAPI.GET("/classes/:classId/overview", partnerHandler.GetClassOverview)

		// 行为记录（老师给学生加成长值的核心接口）
		teacherAPI.POST("/behaviors", behaviorHandler.RecordBehavior)
		teacherAPI.GET("/behaviors", behaviorHandler.GetClassBehaviors)

		// 盲盒管理
		blindboxGroup := teacherAPI.Group("/blindbox")
		{
			blindboxGroup.GET("/pool", blindboxHandler.GetPool)
			blindboxGroup.POST("/pool", blindboxHandler.AddToPool)
			blindboxGroup.DELETE("/pool/:id", blindboxHandler.RemoveFromPool)
			blindboxGroup.POST("/draw", blindboxHandler.DrawForStudent) // 教师代学生开盲盒
		}
	}

	return r
}

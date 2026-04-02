// growth-partner/backend/internal/router/sunshine_router.go
// 阳光章路由

package router

import (
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterSunshineRoutes 注册阳光章路由
func RegisterSunshineRoutes(router *gin.RouterGroup, sunshineHandler *handler.SunshineHandler) {
	// 管理员路由组，需要管理员角色权限
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.RequireRoles("admin"))
	{
		adminRoutes.GET("/sunshine/colors", sunshineHandler.GetSunshineColors)
		adminRoutes.POST("/sunshine/colors", sunshineHandler.CreateSunshineColor)
		adminRoutes.PUT("/sunshine/colors/:id", sunshineHandler.UpdateSunshineColor)
	}

	// 老师路由组，需要老师角色权限
	teacherRoutes := router.Group("/teacher")
	teacherRoutes.Use(middleware.RequireTeacher())
	{
		teacherRoutes.POST("/sunshine/stamp", sunshineHandler.StampSunshine)
		teacherRoutes.GET("/sunshine/stamps", sunshineHandler.GetClassStamps)
		teacherRoutes.POST("/sunshine/awards/evaluate", sunshineHandler.EvaluateSunshineAwards)
		teacherRoutes.GET("/sunshine/awards", sunshineHandler.GetSunshineAwards)
	}

	// 学生路由组，需要学生角色权限
	studentRoutes := router.Group("/student")
	studentRoutes.Use(middleware.RequireStudent())
	{
		studentRoutes.GET("/sunshine/my-stamps", sunshineHandler.GetStudentStamps)
		studentRoutes.GET("/sunshine/my-awards", sunshineHandler.GetStudentAwards)
	}

	// 家长路由组，需要家长角色权限
	parentRoutes := router.Group("/parent")
	parentRoutes.Use(middleware.RequireParent())
	{
		parentRoutes.GET("/children/:childId/sunshine", sunshineHandler.GetChildSunshine)
	}
}

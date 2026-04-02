// growth-partner/backend/internal/router/teacher_router.go
// 老师端模块路由配置

package router

import (
	"growth-partner/config"
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"
	"growth-partner/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// SetupTeacherRoutes 配置老师端相关路由
func SetupTeacherRoutes(
	v1 *gin.RouterGroup,
	cfg *config.Config,
	jwtManager *jwt.Manager,
	teacherHandler *handler.TeacherHandler,
	behaviorHandler *handler.BehaviorHandler,
	broadcastHandler *handler.BroadcastHandler,
	blindboxHandler *handler.BlindboxHandler,
) {
	// ─── 老师端接口（需要 teacher 角色）──────────────────────────
	teacherAPI := v1.Group("/teacher",
		middleware.Auth(jwtManager),
		middleware.RequireTeacher(),
	)
	{
		// 3.1 我的班级管理
		teacherAPI.GET("/my-classes", teacherHandler.GetMyClasses)
		teacherAPI.GET("/classes/:classId/overview", teacherHandler.GetClassOverview)
		teacherAPI.GET("/classes/:classId/students", teacherHandler.GetClassStudents)

		// 3.2 正向行为打分
		teacherAPI.POST("/behaviors", behaviorHandler.RecordBehavior)
		teacherAPI.GET("/behaviors", behaviorHandler.GetClassBehaviors)
		teacherAPI.GET("/behaviors/:id", teacherHandler.GetBehavior)
		teacherAPI.DELETE("/behaviors/:id", teacherHandler.DeleteBehavior)
		teacherAPI.POST("/behaviors/batch", teacherHandler.BatchRecordBehaviors)

		// 3.3 广播发送
		teacherAPI.GET("/broadcasts", teacherHandler.GetBroadcasts)
		teacherAPI.POST("/broadcasts", broadcastHandler.Send)
		teacherAPI.DELETE("/broadcasts/:id", teacherHandler.CancelBroadcast)

		// 3.4 集体挑战管理
		teacherAPI.GET("/challenges", teacherHandler.GetChallenges)
		teacherAPI.POST("/challenges", teacherHandler.CreateChallenge)
		teacherAPI.PATCH("/challenges/:id/complete", teacherHandler.CompleteChallenge)

		// 3.5 题库管理
		teacherAPI.GET("/questions", teacherHandler.GetQuestions)
		teacherAPI.POST("/questions", teacherHandler.CreateQuestion)
		teacherAPI.PUT("/questions/:id", teacherHandler.UpdateQuestion)
		teacherAPI.DELETE("/questions/:id", teacherHandler.DeleteQuestion)
		teacherAPI.POST("/questions/batch-import", teacherHandler.BatchImportQuestions)

		// 3.6 盲盒奖励池管理
		blindboxGroup := teacherAPI.Group("/blindbox")
		{
			blindboxGroup.GET("/pool", blindboxHandler.GetPool)
			blindboxGroup.POST("/pool", blindboxHandler.AddToPool)
			blindboxGroup.PUT("/pool/:id", teacherHandler.UpdateBlindboxPoolItem)
			blindboxGroup.DELETE("/pool/:id", blindboxHandler.RemoveFromPool)
			blindboxGroup.POST("/draw/:childId", blindboxHandler.DrawForStudent)
			blindboxGroup.PATCH("/draws/:drawId/redeem", teacherHandler.ConfirmBlindboxRedeem)
		}

		// 3.7 周报PDF生成
		reportGroup := teacherAPI.Group("/reports")
		{
			reportGroup.POST("/weekly", teacherHandler.GenerateWeeklyReport)
			reportGroup.GET("/weekly", teacherHandler.GetWeeklyReports)
			reportGroup.GET("/weekly/:id/download", teacherHandler.DownloadWeeklyReport)
		}
	}
}

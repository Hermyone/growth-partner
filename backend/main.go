// growth-partner/backend/main.go
// 程序入口：初始化所有依赖并启动 HTTP 服务

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"growth-partner/config"
	"growth-partner/internal/handler"
	"growth-partner/internal/repository"
	"growth-partner/internal/router"
	"growth-partner/internal/service"
	jwtpkg "growth-partner/pkg/jwt"
)

// Version 由构建时 ldflags 注入
var Version = "dev"

func main() {
	log.Printf("🌱 成长伙伴服务启动中，版本: %s", Version)

	// ─── 1. 加载配置 ─────────────────────────────────────────
	cfg := config.Load()

	// ─── 2. 初始化数据库连接 ──────────────────────────────────
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}
	log.Println("✅ PostgreSQL 连接成功")

	// ─── 3. 初始化 Redis ──────────────────────────────────────
	rdb, err := repository.NewRedis(cfg)
	if err != nil {
		log.Fatalf("❌ Redis 连接失败: %v", err)
	}
	log.Println("✅ Redis 连接成功")

	// ─── 4. 初始化 Repository 层 ─────────────────────────────
	userRepo := repository.NewUserRepository(db)
	childRepo := repository.NewChildRepository(db)
	partnerRepo := repository.NewPartnerRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	growthRepo := repository.NewGrowthRepository(db)
	behaviorRepo := repository.NewBehaviorRepository(db)
	battleRepo := repository.NewBattleRepository(db)
	blindRepo := repository.NewBlindboxRepository(db)
	classRepo := repository.NewClassRepository(db)
	schoolRepo := repository.NewSchoolRepository(db)
	adminPermissionRepo := repository.NewAdminPermissionRepository(db)
	parentChildRepo := repository.NewParentChildRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	broadcastRepo := repository.NewBroadcastRepository(db)
	challengeRepo := repository.NewChallengeRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	reportRepo := repository.NewReportRepository(db)
	sunshineRepo := repository.NewSunshineRepository(db)

	milestoneRepo := repository.NewMilestoneRepository(db)
	_ = rdb // Redis 供 Service 层使用

	// ─── 5. 初始化 JWT ────────────────────────────────────────
	jwtManager, _ := jwtpkg.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// ─── 6. 初始化 Service 层 ─────────────────────────────────
	broadcastSvc := service.NewBroadcastService(rdb.Client)
	partnerSvc := service.NewPartnerService(partnerRepo, growthRepo, templateRepo, milestoneRepo, broadcastSvc)
	authSvc := service.NewAuthService(userRepo, childRepo, jwtManager, cfg, rdb)
	behaviorSvc := service.NewBehaviorService(behaviorRepo, partnerSvc, broadcastSvc)
	battleSvc := service.NewBattleService(battleRepo, questionRepo, childRepo)
	blindboxSvc := service.NewBlindboxService(blindRepo, db) // blindboxRepo 待实现
	classSvc := service.NewClassService(classRepo)           // blindboxRepo 待实现
	adminSvc := service.NewAdminService(schoolRepo, classRepo, userRepo, childRepo, parentChildRepo, adminPermissionRepo, auditLogRepo, templateRepo)
	teacherSvc := service.NewTeacherService(classRepo, childRepo, behaviorRepo, partnerRepo, broadcastRepo, challengeRepo, questionRepo, blindRepo, reportRepo, adminPermissionRepo, behaviorSvc, broadcastSvc, blindboxSvc)
	studentSvc := service.NewStudentService(partnerRepo, growthRepo, templateRepo, behaviorRepo, broadcastRepo, milestoneRepo, blindRepo, childRepo)
	parentSvc := service.NewParentService(parentChildRepo, childRepo, partnerRepo, behaviorRepo, broadcastRepo, milestoneRepo, battleRepo)
	templateSvc := service.NewPartnerTemplateService(templateRepo)
	sunshineSvc := service.NewSunshineService(sunshineRepo, classRepo, childRepo, parentChildRepo)

	// ─── 7. 初始化 Handler 层 ─────────────────────────────────
	authHandler := handler.NewAuthHandler(authSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)
	teacherHandler := handler.NewTeacherHandler(teacherSvc)
	studentHandler := handler.NewStudentHandler(studentSvc)
	parentHandler := handler.NewParentHandler(parentSvc)
	partnerHandler := handler.NewPartnerHandler(partnerSvc)
	behaviorHandler := handler.NewBehaviorHandler(behaviorSvc)
	classHandler := handler.NewClassHandler(classSvc) // classRepo 待实现
	childHandler := handler.NewChildHandler(childRepo)
	broadcastHandler := handler.NewBroadcastHandler(broadcastSvc)
	battleHandler := handler.NewBattleHandler(battleSvc)
	blindboxHandler := handler.NewBlindboxHandler(blindboxSvc)
	wsHandler := handler.NewWebSocketHandler()
	templateHandler := handler.NewPartnerTemplateHandler(templateSvc)
	sunshineHandler := handler.NewSunshineHandler(sunshineSvc)

	// ─── 8. 注册路由 ──────────────────────────────────────────
	r := router.SetupRouter(
		cfg, jwtManager,
		authHandler, adminHandler, teacherHandler, studentHandler, parentHandler, partnerHandler, behaviorHandler,
		classHandler, childHandler, broadcastHandler,
		battleHandler, blindboxHandler, wsHandler, templateHandler, sunshineHandler,
	)

	// ─── 9. 启动 HTTP 服务 ────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
	}

	// 在新协程启动服务，主协程监听退出信号
	go func() {
		log.Printf("🚀 服务已启动，监听端口: %d，环境: %s", cfg.App.Port, cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务启动失败: %v", err)
		}
	}()

	// ─── 10. 优雅关闭 ─────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ 正在优雅关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️  服务关闭超时: %v", err)
	}
	log.Println("👋 成长伙伴服务已关闭，再见！")
}

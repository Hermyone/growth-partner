// growth-partner/backend/internal/repository/db.go
// 数据库连接池初始化：使用 GORM 连接 PostgreSQL

package repository

import (
	"growth-partner/config"
	"growth-partner/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB 初始化 GORM 数据库连接
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dbLogger := logger.Default
	if cfg.App.IsDev() {
		dbLogger = dbLogger.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DB.MaxLifetime)

	// 自动迁移表结构（生产环境建议手动管理迁移）
	if cfg.App.IsDev() {
		log.Println("[DB] 正在自动迁移表结构...")
		err = db.AutoMigrate(
			// 账号权限模块
			&model.User{},
			&model.TeacherProfile{},
			&model.AdminPermission{},
			
			// 学籍班级模块
			&model.School{},
			&model.Class{},
			&model.Child{},
			&model.ParentChildRelation{},
			
			// 伙伴系统模块
			&model.PartnerTemplate{},
			&model.Partner{},
			
			// 成长记录模块
			&model.BehaviorRecord{},
			&model.GrowthRecord{},
			&model.Milestone{},
			
			// 知识对战模块
			&model.Question{},
			&model.BattleRecord{},
			
			// 广播与盲盒模块
			&model.BroadcastRecord{},
			&model.BlindBoxPool{},
			&model.BlindBoxDraw{},
		)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

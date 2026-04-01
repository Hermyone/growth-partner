// growth-partner/backend/config/config.go
// 全局配置加载：从环境变量读取，支持 .env 文件，Mac/Windows 均可用

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Config 全局配置结构体
type Config struct {
	App     AppConfig
	DB      DBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Encrypt EncryptConfig
}

// AppConfig 应用基础配置
type AppConfig struct {
	Env          string // "development" | "production"
	Port         int
	Name         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DBConfig PostgreSQL 配置
type DBConfig struct {
	Host         string
	Port         int
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration // 访问令牌有效期
	RefreshTokenTTL time.Duration // 刷新令牌有效期
}

// EncryptConfig 数据加密配置
type EncryptConfig struct {
	AESKey string // 必须是 32 字节（AES-256）
}

// global 全局配置单例
var global *Config

// Load 从环境变量加载配置（Docker 容器内通过 env 注入）
// 不依赖配置文件，完全通过环境变量，Mac/Windows 行为完全一致
func Load() *Config {
	cfg := &Config{
		App: AppConfig{
			Env:          getEnv("APP_ENV", "development"),
			Port:         getEnvInt("APP_PORT", 8080),
			Name:         "成长伙伴",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		DB: DBConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvInt("DB_PORT", 5432),
			Name:         getEnv("DB_NAME", "growth_partner"),
			User:         getEnv("DB_USER", "gp_user"),
			Password:     getEnv("DB_PASSWORD", ""),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxLifetime:  5 * time.Minute,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: 10,
		},
		JWT: JWTConfig{
			Secret:          mustGetEnv("APP_JWT_SECRET"),
			AccessTokenTTL:  24 * time.Hour,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
		Encrypt: EncryptConfig{
			AESKey: mustGetEnv("APP_AES_KEY"),
		},
	}

	// 验证 AES Key 长度
	if len(cfg.Encrypt.AESKey) != 32 {
		log.Fatalf("[Config] AES_KEY 必须恰好是 32 字节，当前长度: %d", len(cfg.Encrypt.AESKey))
	}

	global = cfg
	log.Printf("[Config] 配置加载成功，环境: %s，端口: %d", cfg.App.Env, cfg.App.Port)
	return cfg
}

// Get 获取全局配置（懒加载）
func Get() *Config {
	if global == nil {
		return Load()
	}
	return global
}

// DSN 生成 PostgreSQL 数据源字符串
func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=Asia/Shanghai",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SSLMode,
	)
}

// Addr 生成 Redis 地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDev 是否为开发环境
func (c *AppConfig) IsDev() bool {
	return c.Env == "development"
}

// ─── 内部辅助函数 ──────────────────────────────────────────────

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// mustGetEnv 必须存在的环境变量，缺失则启动失败
func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[Config] 必须设置环境变量: %s", key)
	}
	return v
}

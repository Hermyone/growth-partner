// growth-partner/backend/internal/repository/redis.go
// Redis 连接初始化：用于实时广播和对战状态管理

package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"growth-partner/config"
)

// ─── 常量定义 ──────────────────────────────────────────────────

const (
	// Redis Key 前缀
	KeyPrefixBattle = "battle:" // 对战相关
	KeyPrefixRoom   = "room:"   // 房间相关
	KeyPrefixUser   = "user:"   // 用户相关
	KeyPrefixClass  = "class:"  // 班级相关
	KeyPrefixLock   = "lock:"   // 分布式锁

	// 默认配置
	DefaultPoolSize     = 10
	DefaultMinIdleConns = 5
	DefaultMaxRetries   = 3
)

// ─── Redis 客户端封装 ──────────────────────────────────────────

// RedisClient Redis 客户端封装，扩展原生功能
type RedisClient struct {
	*redis.Client
	config *RedisConfig
}

// RedisConfig Redis 配置（扩展）
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	IdleTimeout  time.Duration
	MaxConnAge   time.Duration
}

// ─── 初始化函数 ──────────────────────────────────────────────────

// NewRedis 初始化 Redis 客户端
func NewRedis(cfg *config.Config) (*RedisClient, error) {
	return NewRedisWithConfig(&RedisConfig{
		Addr:         cfg.Redis.Addr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     getPoolSize(cfg.Redis.PoolSize),
		MinIdleConns: DefaultMinIdleConns,
		MaxRetries:   DefaultMaxRetries,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		IdleTimeout:  5 * time.Minute,
		MaxConnAge:   30 * time.Minute,
	})
}

// NewRedisWithConfig 使用自定义配置初始化 Redis 客户端
func NewRedisWithConfig(cfg *RedisConfig) (*RedisClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis 配置不能为空")
	}

	// 设置默认值
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = DefaultPoolSize
	}
	if cfg.MinIdleConns <= 0 {
		cfg.MinIdleConns = DefaultMinIdleConns
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = DefaultMaxRetries
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 5 * time.Second
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 3 * time.Second
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 3 * time.Second
	}
	if cfg.PoolTimeout <= 0 {
		cfg.PoolTimeout = 4 * time.Second
	}
	if cfg.IdleTimeout <= 0 {
		cfg.IdleTimeout = 5 * time.Minute
	}
	if cfg.MaxConnAge <= 0 {
		cfg.MaxConnAge = 30 * time.Minute
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxRetries:      cfg.MaxRetries,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.IdleTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis 连接失败: %w", err)
	}

	return &RedisClient{
		Client: rdb,
		config: cfg,
	}, nil
}

// ─── 辅助函数 ──────────────────────────────────────────────────

// getPoolSize 获取连接池大小
func getPoolSize(size int) int {
	if size <= 0 {
		return DefaultPoolSize
	}
	return size
}

// ─── RedisClient 扩展方法 ──────────────────────────────────────

// Close 关闭 Redis 连接
func (r *RedisClient) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

// GetConfig 获取 Redis 配置
func (r *RedisClient) GetConfig() *RedisConfig {
	return r.config
}

// Health 健康检查
func (r *RedisClient) Health(ctx context.Context) error {
	return r.Ping(ctx).Err()
}

// ─── Key 生成器（统一管理 Redis Key）──────────────────────────

// BattleKey 生成对战相关 Key
func (r *RedisClient) BattleKey(battleID uint64) string {
	return fmt.Sprintf("%s%d", KeyPrefixBattle, battleID)
}

// BattleUserKey 生成对战用户 Key
func (r *RedisClient) BattleUserKey(battleID, userID uint64) string {
	return fmt.Sprintf("%s%d:user:%d", KeyPrefixBattle, battleID, userID)
}

// RoomKey 生成房间 Key
func (r *RedisClient) RoomKey(roomID string) string {
	return fmt.Sprintf("%s%s", KeyPrefixRoom, roomID)
}

// UserKey 生成用户 Key
func (r *RedisClient) UserKey(userID uint64) string {
	return fmt.Sprintf("%s%d", KeyPrefixUser, userID)
}

// UserBattleKey 生成用户当前对战 Key
func (r *RedisClient) UserBattleKey(userID uint64) string {
	return fmt.Sprintf("%s%d:battle", KeyPrefixUser, userID)
}

// ClassKey 生成班级 Key
func (r *RedisClient) ClassKey(classID uint64) string {
	return fmt.Sprintf("%s%d", KeyPrefixClass, classID)
}

// ClassOnlineUsersKey 生成班级在线用户 Key
func (r *RedisClient) ClassOnlineUsersKey(classID uint64) string {
	return fmt.Sprintf("%s%d:online", KeyPrefixClass, classID)
}

// LockKey 生成分布式锁 Key
func (r *RedisClient) LockKey(lockName string) string {
	return fmt.Sprintf("%s%s", KeyPrefixLock, lockName)
}

// ─── 高级操作方法 ──────────────────────────────────────────────

// SetWithExpire 设置键值并指定过期时间
func (r *RedisClient) SetWithExpire(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// GetOrSet 获取键值，如果不存在则设置
func (r *RedisClient) GetOrSet(ctx context.Context, key string, fn func() (interface{}, error), expiration time.Duration) (interface{}, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err == nil {
		return val, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	// Key 不存在，执行函数获取值
	newVal, err := fn()
	if err != nil {
		return nil, err
	}

	// 设置值
	if err := r.Client.Set(ctx, key, newVal, expiration).Err(); err != nil {
		return nil, err
	}

	return newVal, nil
}

// AcquireLock 获取分布式锁
func (r *RedisClient) AcquireLock(ctx context.Context, lockName string, ttl time.Duration) (bool, error) {
	key := r.LockKey(lockName)

	// 使用 SETNX 实现分布式锁
	ok, err := r.Client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

// ReleaseLock 释放分布式锁
func (r *RedisClient) ReleaseLock(ctx context.Context, lockName string) error {
	key := r.LockKey(lockName)
	return r.Client.Del(ctx, key).Err()
}

// PublishMessage 发布消息（用于实时广播）
func (r *RedisClient) PublishMessage(ctx context.Context, channel string, message interface{}) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

// SubscribeChannel 订阅频道
func (r *RedisClient) SubscribeChannel(ctx context.Context, channels ...string) *redis.PubSub {
	return r.Client.Subscribe(ctx, channels...)
}

// BatchSet 批量设置键值（使用 Pipeline）
func (r *RedisClient) BatchSet(ctx context.Context, pairs map[string]interface{}, expiration time.Duration) error {
	pipe := r.Client.Pipeline()

	for key, value := range pairs {
		pipe.Set(ctx, key, value, expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// IncrementWithExpire 递增并设置过期时间（如果键不存在）
func (r *RedisClient) IncrementWithExpire(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	// 使用 Lua 脚本保证原子性
	script := `
		local current = redis.call('INCR', KEYS[1])
		if current == 1 then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return current
	`

	result, err := r.Client.Eval(ctx, script, []string{key}, expiration.Seconds()).Int64()
	if err != nil {
		return 0, err
	}

	return result, nil
}

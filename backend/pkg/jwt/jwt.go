// growth-partner/backend/pkg/jwt/jwt.go
// JWT 工具：生成/验证 Access Token 和 Refresh Token
// 支持三种角色：学生、老师/园长、家长

package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// ─── 错误定义 ──────────────────────────────────────────────────

var (
	ErrTokenExpired   = errors.New("token 已过期")
	ErrTokenInvalid   = errors.New("token 无效")
	ErrTokenNotBefore = errors.New("token 尚未生效")
	ErrInvalidSecret  = errors.New("无效的密钥配置")
)

// ─── 常量定义 ──────────────────────────────────────────────────

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleParent  = "parent"
	RoleAdmin   = "admin"

	Issuer = "growth-partner"
)

// ─── Claims 定义 ───────────────────────────────────────────────

// Claims JWT 载荷（不含敏感信息，仅用于鉴权）
type Claims struct {
	UserID    uint64 `json:"uid"`
	Username  string `json:"usr"`
	Role      string `json:"rol"`           // "student"/"teacher"/"parent"/"admin"
	ClassID   uint64 `json:"cid,omitempty"` // 班级ID（学生和老师有值）
	ChildID   uint64 `json:"kid,omitempty"` // 学生档案ID（学生角色有值）
	TokenType string `json:"typ"`           // "access"/"refresh"
	jwt.RegisteredClaims
}

// ─── Manager ──────────────────────────────────────────────────

// Manager JWT 管理器
type Manager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewManager 创建 JWT 管理器
func NewManager(secret string, accessTTL, refreshTTL time.Duration) (*Manager, error) {
	if secret == "" {
		return nil, ErrInvalidSecret
	}
	if accessTTL <= 0 || refreshTTL <= 0 {
		return nil, errors.New("token TTL 必须大于 0")
	}

	return &Manager{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}, nil
}

// MustNewManager 创建 JWT 管理器，失败时 panic（用于初始化阶段）
func MustNewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	m, err := NewManager(secret, accessTTL, refreshTTL)
	if err != nil {
		panic(err)
	}
	return m
}

// GenerateAccessToken 生成访问令牌
func (m *Manager) GenerateAccessToken(userID uint64, username, role string, classID, childID uint64) (string, error) {
	return m.generate(userID, username, role, classID, childID, TokenTypeAccess, m.accessTokenTTL)
}

// GenerateRefreshToken 生成刷新令牌
func (m *Manager) GenerateRefreshToken(userID uint64, username, role string, classID, childID uint64) (string, error) {
	return m.generate(userID, username, role, classID, childID, TokenTypeRefresh, m.refreshTokenTTL)
}

// GenerateTokenPair 生成令牌对（同时生成 access 和 refresh token）
func (m *Manager) GenerateTokenPair(userID uint64, username, role string, classID, childID uint64) (accessToken, refreshToken string, err error) {
	accessToken, err = m.GenerateAccessToken(userID, username, role, classID, childID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.GenerateRefreshToken(userID, username, role, classID, childID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generate 内部生成方法
func (m *Manager) generate(userID uint64, username, role string, classID, childID uint64, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		ClassID:   classID,
		ChildID:   childID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    Issuer,
			Subject:   username,
			ID:        generateTokenID(), // 可选：添加唯一 ID 用于追踪
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseAccessToken 解析并验证 Access Token
func (m *Manager) ParseAccessToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, TokenTypeAccess)
}

// ParseRefreshToken 解析并验证 Refresh Token
func (m *Manager) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, TokenTypeRefresh)
}

// parse 内部解析方法
func (m *Manager) parse(tokenStr, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法，防止 alg:none 攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotBefore
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// 验证 Token 类型（防止用 refresh token 访问 access 接口）
	if claims.TokenType != expectedType {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ValidateRole 验证 token 中的角色是否符合要求
func (c *Claims) ValidateRole(allowedRoles ...string) bool {
	if len(allowedRoles) == 0 {
		return true
	}
	for _, role := range allowedRoles {
		if c.Role == role {
			return true
		}
	}
	return false
}

// IsStudent 判断是否为学生角色
func (c *Claims) IsStudent() bool {
	return c.Role == RoleStudent
}

// IsTeacher 判断是否为老师角色
func (c *Claims) IsTeacher() bool {
	return c.Role == RoleTeacher
}

// IsParent 判断是否为家长角色
func (c *Claims) IsParent() bool {
	return c.Role == RoleParent
}

// IsAdmin 判断是否为管理员角色
func (c *Claims) IsAdmin() bool {
	return c.Role == RoleAdmin
}

// HasClassAccess 判断是否有班级访问权限（学生或老师）
func (c *Claims) HasClassAccess() bool {
	return c.IsStudent() || c.IsTeacher()
}

// GetTokenTTL 获取 token 剩余有效时间
func (m *Manager) GetTokenTTL(tokenStr string) (time.Duration, error) {
	claims, err := m.parse(tokenStr, TokenTypeAccess)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, ErrTokenInvalid
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining < 0 {
		return 0, ErrTokenExpired
	}

	return remaining, nil
}

// RefreshAccessToken 使用 refresh token 刷新 access token
func (m *Manager) RefreshAccessToken(refreshTokenStr string) (newAccessToken string, err error) {
	claims, err := m.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		return "", err
	}

	// 生成新的 access token，保持原有信息不变
	return m.GenerateAccessToken(
		claims.UserID,
		claims.Username,
		claims.Role,
		claims.ClassID,
		claims.ChildID,
	)
}

// ─── 辅助函数 ──────────────────────────────────────────────────

// generateTokenID 生成简单的 token ID（可用于日志追踪）
func generateTokenID() string {
	// 可以使用 UUID 或雪花算法，这里简化处理
	return time.Now().Format("20060102150405") + randomString(8)
}

// randomString 生成随机字符串（简化实现，生产环境建议使用 crypto/rand）
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

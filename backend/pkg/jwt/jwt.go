// growth-partner/backend/pkg/jwt/jwt.go
// JWT 工具：生成/验证 Access Token 和 Refresh Token
// 支持三种角色：学生、老师/园长、家长

package jwt

import (
	"errors"
	"time"
)

// ─── 错误定义 ──────────────────────────────────────────────────

var (
	ErrTokenExpired   = errors.New("token 已过期")
	ErrTokenInvalid   = errors.New("token 无效")
	ErrTokenNotBefore = errors.New("token 尚未生效")
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
func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// GenerateAccessToken 生成访问令牌
func (m *Manager) GenerateAccessToken(userID uint64, username, role string, classID, childID uint64) (string, error) {
	return m.generate(userID, username, role, classID, childID, "access", m.accessTokenTTL)
}

// GenerateRefreshToken 生成刷新令牌
func (m *Manager) GenerateRefreshToken(userID uint64, username, role string, classID, childID uint64) (string, error) {
	return m.generate(userID, username, role, classID, childID, "refresh", m.refreshTokenTTL)
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
			Issuer:    "growth-partner",
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseAccessToken 解析并验证 Access Token
func (m *Manager) ParseAccessToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, "access")
}

// ParseRefreshToken 解析并验证 Refresh Token
func (m *Manager) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, "refresh")
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

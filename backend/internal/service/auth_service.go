// growth-partner/backend/internal/service/auth_service.go

package service

import (
	"context"
	"errors"
	"fmt"
	"growth-partner/config"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
	"growth-partner/pkg/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
	ErrRoleMismatch    = errors.New("用户角色不匹配")
)

type AuthService interface {
	Login(ctx context.Context, username, password, role string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetCurrentUser(ctx context.Context, userID uint64) (*model.User, error)
	ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *model.User  `json:"user"`
	Child        *model.Child `json:"child,omitempty"` // 如果是学生，返回学生档案
}

type authServiceImpl struct {
	userRepo    repository.UserRepository
	childRepo   repository.ChildRepository
	jwtManager  *jwt.Manager
	cfg         *config.Config
	redisClient *repository.RedisClient
}

func NewAuthService(u repository.UserRepository, c repository.ChildRepository, j *jwt.Manager, cfg *config.Config, r *repository.RedisClient) AuthService {
	return &authServiceImpl{userRepo: u, childRepo: c, jwtManager: j, cfg: cfg, redisClient: r}
}

func (s *authServiceImpl) Login(ctx context.Context, username, password, role string) (*LoginResponse, error) {
	// 1. 查询用户
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 2. 校验角色
	if string(user.Role) != role {
		return nil, ErrRoleMismatch
	}

	// 3. 校验密码 (使用 bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// 4. 如果是学生角色，获取关联的 ChildID 和 ClassID
	var classID, childID uint64
	var childInfo *model.Child
	if user.Role == model.RoleStudent {
		child, err := s.childRepo.FindByUserID(ctx, user.ID)
		if err == nil && child != nil {
			childID = child.ID
			classID = child.ClassID
			childInfo = child
		}
	}

	// 5. 生成 Tokens
	access, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, string(user.Role), classID, childID)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, string(user.Role), classID, childID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User:         user,
		Child:        childInfo,
	}, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 生成新的 Access Token
	access, err := s.jwtManager.GenerateAccessToken(claims.UserID, claims.Username, claims.Role, claims.ClassID, claims.ChildID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  access,
		RefreshToken: refreshToken, // 沿用旧的 Refresh Token
	}, nil
}

// Logout 用户登出
func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	// 解析 Refresh Token 以获取过期时间
	claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// 计算剩余有效期
	if claims.ExpiresAt == nil {
		return nil // Token 没有过期时间，无需加入黑名单
	}
	duration := time.Until(claims.ExpiresAt.Time)
	if duration <= 0 {
		return nil // Token 已过期，无需加入黑名单
	}

	// 将 Refresh Token 加入 Redis 黑名单
	key := fmt.Sprintf("blacklist:refresh:%s", refreshToken)
	return s.redisClient.SetWithExpire(ctx, key, "1", duration)
}

// GetCurrentUser 获取当前登录用户信息
func (s *authServiceImpl) GetCurrentUser(ctx context.Context, userID uint64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// ChangePassword 修改密码
func (s *authServiceImpl) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	// 1. 查询用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 2. 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// 3. 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. 更新密码
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

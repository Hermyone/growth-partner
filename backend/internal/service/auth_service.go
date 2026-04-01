// growth-partner/backend/internal/service/auth_service.go

package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"growth-partner/config"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
	"growth-partner/pkg/jwt"
)

var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
	ErrRoleMismatch    = errors.New("用户角色不匹配")
)

type AuthService interface {
	Login(ctx context.Context, username, password, role string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *model.User  `json:"user"`
	Child        *model.Child `json:"child,omitempty"` // 如果是学生，返回学生档案
}

type authServiceImpl struct {
	userRepo   repository.UserRepository
	childRepo  repository.ChildRepository
	jwtManager *jwt.Manager
	cfg        *config.Config
}

func NewAuthService(u repository.UserRepository, c repository.ChildRepository, j *jwt.Manager, cfg *config.Config) AuthService {
	return &authServiceImpl{userRepo: u, childRepo: c, jwtManager: j, cfg: cfg}
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

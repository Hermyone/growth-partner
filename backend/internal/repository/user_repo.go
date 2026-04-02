// growth-partner/backend/internal/repository/user_repo.go

package repository

import (
	"context"
	"growth-partner/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

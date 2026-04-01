// growth-partner/backend/internal/service/class_service.go
// 班级服务：管理班级生命周期及班级信息

package service

import (
	"context"
	"fmt"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// ClassService 班级服务接口
type ClassService interface {
	// 创建班级
	CreateClass(ctx context.Context, name string, grade int, schoolYear string, teacherID uint64) (*model.Class, error)
	// 根据ID获取班级
	GetClassByID(ctx context.Context, classID uint64) (*model.Class, error)
	// 获取教师负责的所有活跃班级
	GetClassesByTeacher(ctx context.Context, teacherID uint64) ([]*model.Class, error)
}

type classServiceImpl struct {
	classRepo repository.ClassRepository
}

// NewClassService 创建班级服务实例
func NewClassService(repo repository.ClassRepository) ClassService {
	return &classServiceImpl{classRepo: repo}
}

// CreateClass 创建一个新的班级
func (s *classServiceImpl) CreateClass(ctx context.Context, name string, grade int, schoolYear string, teacherID uint64) (*model.Class, error) {
	class := &model.Class{
		Name:       name,
		Grade:      grade,
		SchoolYear: schoolYear,
		TeacherID:  teacherID,
		IsActive:   true,
	}

	if err := s.classRepo.Create(ctx, class); err != nil {
		return nil, fmt.Errorf("数据库创建班级失败: %w", err)
	}

	return class, nil
}

// GetClassByID 获取指定班级的详细信息
func (s *classServiceImpl) GetClassByID(ctx context.Context, classID uint64) (*model.Class, error) {
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("查询班级失败: %w", err)
	}
	return class, nil
}

// GetClassesByTeacher 获取教师管理的所有班级列表
func (s *classServiceImpl) GetClassesByTeacher(ctx context.Context, teacherID uint64) ([]*model.Class, error) {
	classes, err := s.classRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("获取教师班级列表失败: %w", err)
	}
	return classes, nil
}

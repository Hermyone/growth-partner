// growth-partner/backend/internal/repository/question_repo.go
// 题库仓库接口和实现

package repository

import (
	"context"
	"growth-partner/internal/model"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *model.Question) error
	FindByID(ctx context.Context, id uint64) (*model.Question, error)
	FindAll(ctx context.Context, params map[string]interface{}) ([]*model.Question, int64, error)
	FindByClassID(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Question, int64, error)
	Update(ctx context.Context, question *model.Question) error
	Delete(ctx context.Context, id uint64) error
	BatchCreate(ctx context.Context, questions []*model.Question) error
}

type questionRepositoryImpl struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepositoryImpl{db: db}
}

func (r *questionRepositoryImpl) Create(ctx context.Context, question *model.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.Question, error) {
	var question model.Question
	err := r.db.WithContext(ctx).First(&question, id).Error
	return &question, err
}

func (r *questionRepositoryImpl) FindAll(ctx context.Context, params map[string]interface{}) ([]*model.Question, int64, error) {
	var questions []*model.Question
	db := r.db.WithContext(ctx).Where("is_active = ?", true)

	// 应用过滤条件
	if subjectID, ok := params["subject_id"].(uint64); ok && subjectID > 0 {
		db = db.Where("subject_id = ?", subjectID)
	}
	if questionType, ok := params["question_type"].(string); ok && questionType != "" {
		db = db.Where("question_type = ?", questionType)
	}
	if difficulty, ok := params["difficulty"].(int); ok && difficulty > 0 {
		db = db.Where("difficulty = ?", difficulty)
	}
	if isPublic, ok := params["is_public"].(bool); ok {
		db = db.Where("is_public = ?", isPublic)
	}
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		db = db.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 分页
	var count int64
	if err := db.Model(&model.Question{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "created_at DESC"
	if o, ok := params["order"].(string); ok && o != "" {
		order = o
	}
	db = db.Order(order)

	// 分页
	if page, ok := params["page"].(int); ok && page > 0 {
		limit := 10
		if l, ok := params["limit"].(int); ok && l > 0 {
			limit = l
		}
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	err := db.Find(&questions).Error
	return questions, count, err
}

func (r *questionRepositoryImpl) FindByClassID(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Question, int64, error) {
	var questions []*model.Question
	db := r.db.WithContext(ctx).Where("(class_id = ? OR is_public = ?) AND is_active = ?", classID, true, true)

	// 应用过滤条件
	if subjectID, ok := params["subject_id"].(uint64); ok && subjectID > 0 {
		db = db.Where("subject_id = ?", subjectID)
	}
	if questionType, ok := params["question_type"].(string); ok && questionType != "" {
		db = db.Where("question_type = ?", questionType)
	}
	if difficulty, ok := params["difficulty"].(int); ok && difficulty > 0 {
		db = db.Where("difficulty = ?", difficulty)
	}
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		db = db.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 分页
	var count int64
	if err := db.Model(&model.Question{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "created_at DESC"
	if o, ok := params["order"].(string); ok && o != "" {
		order = o
	}
	db = db.Order(order)

	// 分页
	if page, ok := params["page"].(int); ok && page > 0 {
		limit := 10
		if l, ok := params["limit"].(int); ok && l > 0 {
			limit = l
		}
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	err := db.Find(&questions).Error
	return questions, count, err
}

func (r *questionRepositoryImpl) Update(ctx context.Context, question *model.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

func (r *questionRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Question{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *questionRepositoryImpl) BatchCreate(ctx context.Context, questions []*model.Question) error {
	return r.db.WithContext(ctx).Create(&questions).Error
}

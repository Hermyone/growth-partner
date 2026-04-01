// growth-partner/backend/internal/model/question.go
// 题库模型：支持多科目，老师可自定义上传

package model

// QuestionType 题目类型
type QuestionType string

const (
	QTypeChoice QuestionType = "choice" // 选择题
	QTypeFill   QuestionType = "fill"   // 填空题
	QTypePoetry QuestionType = "poetry" // 诗词接龙（特殊类型）
	QTypeMath   QuestionType = "math"   // 速算题（程序自动生成）
)

// Question 题库表
type Question struct {
	Base

	ClassID uint64        `gorm:"index"           json:"class_id"` // 0=系统题库，>0=班级专属
	Subject BattleSubject `gorm:"size:16;index;not null" json:"subject"`
	Type    QuestionType  `gorm:"size:16;not null"       json:"type"`
	Grade   int           `gorm:"index"                   json:"grade"` // 适合年级

	Content string `gorm:"size:512;not null" json:"content"`           // 题目内容
	Options JSON   `gorm:"type:jsonb"        json:"options"`           // 选项（选择题用）
	Answer  string `gorm:"size:256;not null" json:"-"`                 // 答案（不暴露给学生）
	Explain string `gorm:"size:512"          json:"explain,omitempty"` // 解析

	Difficulty int    `gorm:"default:1;index" json:"difficulty"` // 1-3 难度
	IsActive   bool   `gorm:"default:true"    json:"is_active"`
	CreatedBy  uint64 `gorm:"index"         json:"created_by"` // 创建者（老师ID）
}

func (Question) TableName() string { return "questions" }

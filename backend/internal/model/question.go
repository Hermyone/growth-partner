// growth-partner/backend/internal/model/question.go
// 题库模型

package model

// QuestionType 题目类型
type QuestionType string

const (
	QuestionTypeSingle  QuestionType = "single"  // 单选题
	QuestionTypeMultiple QuestionType = "multiple" // 多选题
	QuestionTypeJudge   QuestionType = "judge"   // 判断题
	QuestionTypeFill    QuestionType = "fill"    // 填空题
)

// Question 题库表
type Question struct {
	Base

	SubjectID      uint64        `gorm:"index;not null" json:"subject_id"`      // 科目ID
	Content        string        `gorm:"type:text;not null" json:"content"`        // 题目内容
	QuestionType   QuestionType  `gorm:"size:20;not null" json:"question_type"`   // 题目类型
	Options        string        `gorm:"type:text" json:"options,omitempty"`       // 选项（JSON格式存储）
	Answer         string        `gorm:"type:text;not null" json:"answer"`         // 答案
	Explanation    string        `gorm:"type:text" json:"explanation,omitempty"`    // 解析
	Difficulty     int           `gorm:"default:1;check:difficulty >= 1 AND difficulty <= 5" json:"difficulty"` // 难度等级（1-5）
	IsPublic       bool          `gorm:"default:false" json:"is_public"`          // 是否为公共题目
	ClassID        *uint64       `gorm:"index" json:"class_id,omitempty"`        // 所属班级ID（非公共题目时必填）
	CreatedBy      uint64        `gorm:"index;not null" json:"created_by"`        // 创建人ID（老师）
	IsActive       bool          `gorm:"default:true" json:"is_active"`           // 是否启用
}

func (Question) TableName() string { return "questions" }

package system

import "time"

type ExamPaperAnswer_1 struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement"`
	ExamPaperID     int       `gorm:"column:exam_paper_id"`
	PaperName       string    `gorm:"column:paper_name;type:varchar(255)"`
	PaperType       int       `gorm:"column:paper_type"`
	SubjectID       int       `gorm:"column:subject_id"`
	SystemScore     int       `gorm:"column:system_score"`
	UserScore       int       `gorm:"column:user_score"`
	PaperScore      int       `gorm:"column:paper_score"`
	QuestionCorrect int       `gorm:"column:question_correct"`
	QuestionCount   int       `gorm:"column:question_count"`
	DoTime          int       `gorm:"column:do_time"`
	Status          int       `gorm:"column:status"`
	CreateUser      int       `gorm:"column:create_user"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime"`
	TaskExamID      int       `gorm:"column:task_exam_id"`
}

// TableName 指定表名
func (ExamPaperAnswer_1) TableName() string {
	return "t_exam_paper_answer"
}

type ExamPaperQuestionCustomerAnswer_1 struct {
	ID                    int `gorm:"primaryKey"`
	QuestionID            int
	ExamPaperID           int
	ExamPaperAnswerID     int
	QuestionType          int
	SubjectID             int
	CustomerScore         int
	QuestionScore         int
	QuestionTextContentID int
	Answer                *string
	TextContentID         int
	DoRight               bool
	CreateUser            int
	CreateTime            time.Time
	ItemOrder             int
}

func (ExamPaperQuestionCustomerAnswer_1) TableName() string {
	return "t_exam_paper_question_customer_answer"
}

type Question_1 struct {
	ID                int       `gorm:"column:id;primaryKey;autoIncrement"`
	QuestionType      int       `gorm:"column:question_type"`
	SubjectID         int       `gorm:"column:subject_id"`
	Score             int       `gorm:"column:score"`
	GradeLevel        int       `gorm:"column:grade_level"`
	Difficult         int       `gorm:"column:difficult"`
	Correct           string    `gorm:"column:correct;type:text"`
	InfoTextContentID int       `gorm:"column:info_text_content_id"`
	CreateUser        int       `gorm:"column:create_user"`
	Status            int       `gorm:"column:status"`
	CreateTime        time.Time `gorm:"column:create_time;type:datetime"`
	Deleted           []byte    `gorm:"column:deleted"`
}

func (Question_1) TableName() string {
	return "t_question"
}

type QuestionErrorCount struct {
	ID            int       `gorm:"primaryKey"`
	QuestionID    int       `gorm:"not null"`
	ErrorCount    int       `gorm:"not null;default:0"`
	LastErrorTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (QuestionErrorCount) TableName() string {
	return "t_question_error_count"
}

type UserWrongBook struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	UserID      int       `gorm:"column:user_id;not null"`
	QuestionID  int       `gorm:"column:question_id;not null"`
	ExamPaperID int       `gorm:"column:exam_paper_id;not null"`
	SubjectID   int       `gorm:"column:subject_id;not null"` // 新增字段
	CreateTime  time.Time `gorm:"column:create_time;not null"`
	UpdateTime  time.Time `gorm:"column:update_time;not null"`
	ErrorCount  int       `gorm:"column:error_count;not null;default:1"`
}

func (UserWrongBook) TableName() string {
	return "t_user_wrong_book"
}

type ExamPaperTextContent1 struct {
	ID         int        `gorm:"column:id"`
	Content    string     `gorm:"column:content"`
	CreateTime *time.Time `gorm:"column:create_time"`
}

func (ExamPaperTextContent1) TableName() string {
	return "t_text_content"
}

// PaperStruct 用于解析试卷内容的结构体
type PaperStructWrapper []PaperStruct

type PaperStruct struct {
	Name          string          `json:"name"`
	QuestionItems []QuestionItem1 `json:"questionItems"`
}

type QuestionItem1 struct {
	ID        int `json:"id"`
	ItemOrder int `json:"itemOrder"`
}

// ExamPaper 结构体
type ExamPaper_1 struct {
	ID                 int        `gorm:"primaryKey;autoIncrement"`              // 主键，自增
	Name               string     `gorm:"type:varchar(255);not null"`            // 试卷名称
	SubjectID          int        `gorm:"not null"`                              // 学科
	PaperType          int        `gorm:"not null"`                              // 试卷类型 (1.固定试卷 4.时段试卷 6.任务试卷)
	GradeLevel         int        `gorm:"not null"`                              // 年级
	Score              int        `gorm:"not null"`                              // 试卷总分(千分制)
	QuestionCount      int        `gorm:"not null"`                              // 题目数量
	SuggestTime        int        `gorm:"not null"`                              // 建议时长(分钟)
	LimitStartTime     *time.Time `gorm:"type:datetime"`                         // 时段试卷 开始时间
	LimitEndTime       *time.Time `gorm:"type:datetime"`                         // 时段试卷 结束时间
	FrameTextContentID int        `gorm:"column:frame_text_content_id;not null"` // 试卷框架 内容为JSON
	CreateUser         int        `gorm:"not null"`                              // 创建用户
	CreateTime         time.Time  `gorm:"type:datetime;not null"`                // 创建时间
	Deleted            []byte     `gorm:"not null;default:false"`                // 是否删除
	TaskExamID         *int       `gorm:"default:null"`                          // 任务试卷ID
}

func (ExamPaper_1) TableName() string {
	return "t_exam_paper"
}

type QuestionItem_1 struct {
	ID        int `json:"id"`
	ItemOrder int `json:"itemOrder"`
}

type TitleItem struct {
	Name          string           `json:"name"`
	QuestionItems []QuestionItem_1 `json:"questionItems"`
}

type PaperVisibility struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	PaperID   int       `gorm:"column:paper_id;not null"`
	UserID    int       `gorm:"column:user_id;not null"`
	CreatedBy int       `gorm:"column:created_by;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

// TableName sets the insert table name for this struct type
func (PaperVisibility) TableName() string {
	return "t_paper_visibility"
}

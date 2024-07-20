package system

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/slyrx/gin_exam_system/server/others/global"
)

type SysUser struct {
	global.GES_MODEL
	UUID        uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`                                                     // 用户UUID
	Username    string         `json:"userName" gorm:"index;comment:用户登录名"`                                                  // 用户登录名
	Password    string         `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName    string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email       string         `json:"email"  gorm:"comment:用户邮箱"`                      // 用户邮箱
	Enable      int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

func (SysUser) TableName() string {
	return "sys_users"
}

// TUser represents a user in the t_user table
type SysExamUser struct {
	ID             int       `gorm:"column:id;primary_key;auto_increment"`
	UserUUID       string    `gorm:"column:user_uuid;size:36"`
	UserName       string    `gorm:"column:user_name;size:255"`
	Password       string    `gorm:"column:password;size:255"`
	RealName       string    `gorm:"column:real_name;size:255"`
	Age            int       `gorm:"column:age"`
	Sex            int       `gorm:"column:sex"`
	BirthDay       time.Time `gorm:"column:birth_day"`
	UserLevel      int       `gorm:"column:user_level"`
	Phone          string    `gorm:"column:phone;size:255"`
	Role           int       `gorm:"column:role"`
	Status         int       `gorm:"column:status"`
	ImagePath      string    `gorm:"column:image_path;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	ModifyTime     time.Time `gorm:"column:modify_time"`
	LastActiveTime time.Time `gorm:"column:last_active_time"`
	Deleted        []byte    `gorm:"column:deleted"`
	WxOpenID       string    `gorm:"column:wx_open_id;size:255"`
}

// TableName specifies the table name for TUser struct
func (SysExamUser) TableName() string {
	return "t_user"
}

// ExamPaperSubmitVM 定义了提交试卷的结构体
type ExamPaperSubmitVM struct {
	QuestionID  *int                    `json:"questionId"` // 使用指针以允许 null 值
	DoTime      int                     `json:"doTime"`
	AnswerItems []ExamPaperSubmitItemVM `json:"answerItems"`
	ID          int                     `json:"id"`
}

// ExamPaperSubmitItemVM 定义了每个回答项的结构体
type ExamPaperSubmitItemVM struct {
	QuestionID   int      `json:"questionId"`
	Content      string   `json:"content"`
	ContentArray []string `json:"contentArray"`
	Completed    bool     `json:"completed"`
	ItemOrder    int      `json:"itemOrder"`
}

type ExamPaperAnswerInfo struct {
	ExamPaper                        ExamPaper
	ExamPaperAnswer                  ExamPaperAnswer
	ExamPaperQuestionCustomerAnswers []ExamPaperQuestionCustomerAnswer
}

type ExamPaperAnswer struct {
	ID              int       `gorm:"column:id"`
	UserID          int       `gorm:"column:user_id"`
	ExamPaperID     int       `gorm:"column:exam_paper_id"`
	PaperName       string    `gorm:"column:paper_name"`
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
	CreateTime      time.Time `gorm:"column:create_time"`
	TaskExamID      int
}

// TableName 设置表名
func (ExamPaperAnswer) TableName() string {
	return "t_exam_paper_answer"
}

// ExamPaper 结构体定义
type ExamPaper struct {
	ID                 int       `gorm:"column:id"`
	Name               string    `gorm:"column:name"`
	SubjectID          int       `gorm:"column:subject_id"`
	PaperType          int       `gorm:"column:paper_type"`
	GradeLevel         int       `gorm:"column:grade_level"`
	Score              int       `gorm:"column:score"`
	QuestionCount      int       `gorm:"column:question_count"`
	SuggestTime        int       `gorm:"column:suggest_time"`
	LimitStartTime     time.Time `gorm:"column:limit_start_time"`
	LimitEndTime       time.Time `gorm:"column:limit_end_time"`
	FrameTextContentID int       `gorm:"column:frame_text_content_id"`
	CreateUser         int       `gorm:"column:create_user"`
	CreateTime         time.Time `gorm:"column:create_time"`
	Deleted            []byte    `gorm:"column:deleted"`
	TaskExamID         int       `gorm:"column:task_exam_id"`
}

// TableName 指定表名
func (ExamPaper) TableName() string {
	return "t_exam_paper"
}

type ExamPaperTextContent struct {
	ID         int       `gorm:"column:id"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (ExamPaperTextContent) TableName() string {
	return "t_text_content"
}

type ExamPaperTitleItemObject struct {
	Name          string         `json:"name"`
	QuestionItems []QuestionItem `json:"questionItems"`
}

type QuestionItem struct {
	ID        int `json:"id"`
	ItemOrder int `json:"itemOrder"`
}

type ExamPaperQuestionCustomerAnswer struct {
	QuestionID            int       `gorm:"column:question_id"`
	QuestionScore         int       `gorm:"column:question_score"`
	SubjectID             int       `gorm:"column:subject_id"`
	CreateTime            time.Time `gorm:"column:create_time"`
	CreateUser            int       `gorm:"column:create_user"`
	TextContentID         *int      `gorm:"column:text_content_id"`
	ExamPaperID           int       `gorm:"column:exam_paper_id"`
	QuestionType          int       `gorm:"column:question_type"`
	Answer                *string   `gorm:"column:answer"`
	CustomerScore         int       `gorm:"column:customer_score"`
	ExamPaperAnswerID     int       `gorm:"column:exam_paper_answer_id"`
	DoRight               bool      `gorm:"column:do_right"`
	QuestionTextContentID int       `gorm:"column:question_text_content_id"`
	ItemOrder             int       `gorm:"column:item_order"`
	ErrCount              int       `gorm:"column:err_count"`
}

// TableName 设置表名
func (ExamPaperQuestionCustomerAnswer) TableName() string {
	return "t_exam_paper_question_customer_answer"
}

// Question 结构体
type Question struct {
	ID                 int                 `gorm:"column:id"`
	QuestionType       int                 `gorm:"column:question_type"`
	SubjectID          int                 `gorm:"column:subject_id"`
	Score              int                 `gorm:"column:score"`
	GradeLevel         string              `gorm:"column:grade_level"`
	Difficult          string              `gorm:"column:difficult"`
	Correct            string              `gorm:"column:correct"`
	InfoTextContentID  int                 `gorm:"column:info_text_content_id"`
	CreateUser         string              `gorm:"column:create_user"`
	Status             string              `gorm:"column:status"`
	CreateTime         time.Time           `gorm:"column:create_time"`
	Deleted            []byte              `gorm:"column:deleted"`
	UserQuestionErrors []UserQuestionError `gorm:"foreignKey:QuestionID"`
	ErrCount           int                 `gorm:"-"`
	ErrCountTotal      int                 `gorm:"column:err_count_total"`
}

func (Question) TableName() string {
	return "t_question"
}

type UserEventLog struct {
	UserID     int       `gorm:"column:user_id"`
	UserName   string    `gorm:"column:user_name"`
	RealName   string    `gorm:"column:real_name"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

// TableName 设置表名
func (UserEventLog) TableName() string {
	return "t_user_event_log"
}

type QuestionObject struct {
	TitleContent        string        `json:"titleContent"`
	Analyze             string        `json:"analyze"`
	QuestionItemObjects []interface{} `json:"questionItemObjects"`
	Correct             string        `json:"correct"`
}

// 定义表的结构体
type UserQuestionError struct {
	UserID     uint `gorm:"column:user_id"`
	QuestionID uint `gorm:"column:question_id"`
	ErrCount   uint `gorm:"column:err_count"`
}

func (UserQuestionError) TableName() string {
	return "t_user_question_errors"
}

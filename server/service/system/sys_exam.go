package system

import (
	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ExamService struct{}

var ExamServiceApp = new(ExamService)

func (examService *ExamService) CreateUserEventLog(userEventLog systemMod.UserEventLog) (err error) {
	return global.GES_DB.Debug().Create(&userEventLog).Error
}

func (examService *ExamService) selectTextContentByID(id int) (s string, err error) {
	examPaperTextContent := &systemMod.ExamPaperTextContent{}
	err = global.GES_DB.Debug().Where("id = ?", id).First(examPaperTextContent).Error
	s = examPaperTextContent.Content
	return
}

func (examService *ExamService) selectQuestionsByIds(userId int, ids []int) ([]systemMod.Question, error) {
	var questions []systemMod.Question

	// 查询问题，并预加载 UserQuestionErrors
	err := global.GES_DB.Debug().
		Preload("UserQuestionErrors", func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id = ?", userId)
		}).
		Where("id IN (?)", ids).
		Find(&questions).Error
	if err != nil {
		return nil, err
	}

	// 将 UserQuestionErrors 中的 err_count 赋值给 Question 的 ErrCount
	for i := range questions {
		questions[i].ErrCount = 0
		for _, userError := range questions[i].UserQuestionErrors {
			questions[i].ErrCount += int(userError.ErrCount)
		}
	}

	return questions, nil
}

func (examService *ExamService) QuestionTypeEnumNeedSaveTextContent(questionType int) bool {
	// Implement logic based on your enum
	return false
}

func (examService *ExamService) GetUserInfo(userName string) *systemMod.SysExamUser {
	// Implement logic based on your enum
	var user systemMod.SysExamUser
	// Query the database for a user with the specified username
	err := global.GES_DB.Debug().Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		return nil
	}
	global.GES_LOG.Info("exam3", zap.Any("GetUserInfo", user))
	return &user
}

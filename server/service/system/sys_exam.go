package system

import (
	"encoding/json"
	"time"

	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ExamService struct{}

var ExamServiceApp = new(ExamService)

func (examService *ExamService) CalculateExamPaperAnswer(user *systemMod.SysExamUser, examPaperSubmitVM systemMod.ExamPaperSubmitVM) (examPaperAnswerInfo *systemMod.ExamPaperAnswerInfo, err error) {
	// 创建一个新的ExamPaperAnswerInfo对象
	examPaperAnswerInfo = &systemMod.ExamPaperAnswerInfo{}
	now := time.Now()                                                          // 获取当前时间
	examPaper, err := ExamServiceApp.selectExamPaperByID(examPaperSubmitVM.ID) // 根据提交的试卷ID从数据库中获取ExamPaper对象
	if err != nil {
		return
	}
	// 根据试卷类型码获取ExamPaperTypeEnum枚举类型
	// 任务试卷只能做一次
	global.GES_LOG.Info("exam", zap.Int("DB init", examPaper.PaperType), zap.Int("DB init", examPaperSubmitVM.ID))
	if examPaper.PaperType == 1 {
		examPaperAnswer := ExamServiceApp.getExamPaperAnswerByPidUid(examPaperSubmitVM.ID, 2)
		if examPaperAnswer != nil {
			return nil, err
		}
	}
	global.GES_LOG.Info("exam", zap.Int("2", examPaper.PaperType))
	// 如果是任务类型的试卷，检查用户是否已经提交过该试卷答案，如果提交过则返回null
	// 获取试卷结构的文本内容
	// 将试卷结构的JSON文本内容转换为ExamPaperTitleItemObject对象列表
	// 获取所有问题的ID列表
	// 根据问题ID列表从数据库中获取Question对象列表
	// 将题目结构转换为题目答案
	// 获取当前问题的Question对象
	// 获取用户提交的当前问题的答案
	// 创建ExamPaperQuestionCustomerAnswer对象
	// 创建ExamPaperAnswer对象
	// 设置ExamPaperAnswerInfo的属性
	// 返回ExamPaperAnswerInfo对象
	global.GES_LOG.Info("exam", zap.Int("DB init", examPaper.FrameTextContentID))
	frameTextContent, err := ExamServiceApp.selectTextContentByID(examPaper.FrameTextContentID)
	if err != nil {
		return nil, err
	}
	var examPaperTitleItemObjects []systemMod.ExamPaperTitleItemObject
	json.Unmarshal([]byte(frameTextContent), &examPaperTitleItemObjects)

	var questionIds []int
	for _, item := range examPaperTitleItemObjects {
		for _, q := range item.QuestionItems {
			questionIds = append(questionIds, q.ID)
		}
	}

	global.GES_LOG.Info("exam", zap.Any("questionIds", questionIds))
	questions, err := ExamServiceApp.selectQuestionsByIds(questionIds)
	if err != nil {
		return nil, err
	}
	global.GES_LOG.Info("exam", zap.Any("questions", questions))
	var examPaperQuestionCustomerAnswers []systemMod.ExamPaperQuestionCustomerAnswer
	for _, item := range examPaperTitleItemObjects {
		for _, q := range item.QuestionItems {
			var question systemMod.Question
			for _, qu := range questions {
				if qu.ID == q.ID {
					question = qu
					break
				}
			}

			var customerQuestionAnswer *systemMod.ExamPaperSubmitItemVM
			for _, tq := range examPaperSubmitVM.AnswerItems {
				if tq.QuestionID == q.ID {
					customerQuestionAnswer = &tq
					break
				}
			}
			global.GES_LOG.Info("examPaperTitleItemObjects item", zap.Any("questions", question))
			answer := ExamServiceApp.ExamPaperQuestionCustomerAnswerFromVM(question, customerQuestionAnswer, examPaper, q.ItemOrder, user, now)
			examPaperQuestionCustomerAnswers = append(examPaperQuestionCustomerAnswers, answer)
		}
	}
	global.GES_LOG.Info("exam", zap.Any("questions", examPaperQuestionCustomerAnswers))
	// 计算每个题的对错，计算总得分
	examPaperAnswer := ExamServiceApp.ExamPaperAnswerFromVM(examPaperSubmitVM, examPaper, examPaperQuestionCustomerAnswers, user, now)
	examPaperAnswerInfo.ExamPaper = examPaper
	examPaperAnswerInfo.ExamPaperAnswer = examPaperAnswer
	examPaperAnswerInfo.ExamPaperQuestionCustomerAnswers = examPaperQuestionCustomerAnswers

	// 数据库保存
	err = ExamServiceApp.createExamPaperAnswerFromVM(examPaperAnswer)
	if err != nil {
		return nil, err
	}

	// 更新错题数
	err = ExamServiceApp.createQuestion(questionIds)
	if err != nil {
		return nil, err
	}

	err = ExamServiceApp.createExamPaperQuestionCustomerAnswer(examPaperQuestionCustomerAnswers)
	if err != nil {
		return nil, err
	}
	return examPaperAnswerInfo, err
}

func (examService *ExamService) CreateUserEventLog(userEventLog systemMod.UserEventLog) (err error) {
	return global.GES_DB.Debug().Create(&userEventLog).Error
}

func (apiService *ExamService) createQuestion(questionIds []int) (err error) {
	err = global.GES_DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&systemMod.Question{}).
			Where("id IN ?", questionIds).
			Updates(map[string]interface{}{"err_count": gorm.Expr("err_count + 1")})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected != 2 {
			// return fmt.Errorf("预期更新 2 条记录，实际更新 %d 条", result.RowsAffected)
			return nil
		}

		return nil
	})
	return
}

func (apiService *ExamService) createExamPaperQuestionCustomerAnswer(examPaperQuestionCustomerAnswer []systemMod.ExamPaperQuestionCustomerAnswer) (err error) {
	return global.GES_DB.Debug().Create(&examPaperQuestionCustomerAnswer).Error
}

func (apiService *ExamService) createExamPaperAnswerFromVM(examPaperAnswer systemMod.ExamPaperAnswer) (err error) {
	return global.GES_DB.Debug().Create(&examPaperAnswer).Error
}

func (examService *ExamService) selectExamPaperByID(id int) (examPaper systemMod.ExamPaper, err error) {
	err = global.GES_DB.Debug().Where("id = ?", id).First(&examPaper).Error
	return
}

func (examService *ExamService) getExamPaperAnswerByPidUid(pid, uid int) *systemMod.ExamPaperAnswer {
	// 模拟数据库查询
	return nil
}

func (examService *ExamService) selectTextContentByID(id int) (s string, err error) {
	examPaperTextContent := &systemMod.ExamPaperTextContent{}
	err = global.GES_DB.Debug().Where("id = ?", id).First(examPaperTextContent).Error
	s = examPaperTextContent.Content
	return
}

func (examService *ExamService) selectQuestionsByIds(ids []int) (question []systemMod.Question, err error) {
	err = global.GES_DB.Debug().Where("id in (?)", ids).Find(&question).Error
	if err != nil {
		return
	}
	// 模拟数据库查询
	return
}

func (examService *ExamService) ExamPaperQuestionCustomerAnswerFromVM(question systemMod.Question, customerQuestionAnswer *systemMod.ExamPaperSubmitItemVM, examPaper systemMod.ExamPaper, itemOrder int, user *systemMod.SysExamUser, now time.Time) systemMod.ExamPaperQuestionCustomerAnswer {
	global.GES_LOG.Info("exam0", zap.Any("questions", customerQuestionAnswer))
	// 判断每个题目的正确情况
	examPaperQuestionCustomerAnswer := systemMod.ExamPaperQuestionCustomerAnswer{
		QuestionID:            question.ID,
		ExamPaperID:           examPaper.ID,
		QuestionScore:         question.Score,
		SubjectID:             examPaper.SubjectID,
		ItemOrder:             itemOrder,
		CreateTime:            now,
		CreateUser:            2,
		QuestionType:          question.QuestionType,
		QuestionTextContentID: question.InfoTextContentID,
	}
	global.GES_LOG.Info("exam1", zap.Any("questions", customerQuestionAnswer))
	if customerQuestionAnswer == nil {
		global.GES_LOG.Info("exam1", zap.Any("nil", customerQuestionAnswer))
		examPaperQuestionCustomerAnswer.CustomerScore = 0
	} else {
		global.GES_LOG.Info("exam1", zap.Any("not nil", customerQuestionAnswer))
		ExamServiceApp.setSpecialFromVM(&examPaperQuestionCustomerAnswer, question, *customerQuestionAnswer)
	}
	return examPaperQuestionCustomerAnswer

}

func (examService *ExamService) setSpecialFromVM(examPaperQuestionCustomerAnswer *systemMod.ExamPaperQuestionCustomerAnswer, question systemMod.Question, customerQuestionAnswer systemMod.ExamPaperSubmitItemVM) {
	global.GES_LOG.Info("exam2", zap.Int("setSpecialFromVM", question.QuestionType))
	switch question.QuestionType {
	case global.SingleChoice:
		global.GES_LOG.Info("exam2", zap.Int("setSpecialFromVM", 1))
		examPaperQuestionCustomerAnswer.Answer = &customerQuestionAnswer.Content
		examPaperQuestionCustomerAnswer.DoRight = question.Correct == customerQuestionAnswer.Content
		if examPaperQuestionCustomerAnswer.DoRight {
			examPaperQuestionCustomerAnswer.CustomerScore = question.Score
		} else {
			examPaperQuestionCustomerAnswer.CustomerScore = 0
			examPaperQuestionCustomerAnswer.ErrCount = question.ErrCount + 1
		}
	case global.MultipleChoice:
		customerAnswer := utils.ContentToString(customerQuestionAnswer.ContentArray)
		examPaperQuestionCustomerAnswer.Answer = &customerAnswer
		examPaperQuestionCustomerAnswer.DoRight = customerAnswer == question.Correct
		if examPaperQuestionCustomerAnswer.DoRight {
			examPaperQuestionCustomerAnswer.CustomerScore = question.Score
		} else {
			examPaperQuestionCustomerAnswer.CustomerScore = 0
		}
	case global.GapFilling:
		correctAnswer := utils.ToJsonStr(customerQuestionAnswer.ContentArray)
		examPaperQuestionCustomerAnswer.Answer = &correctAnswer
		examPaperQuestionCustomerAnswer.CustomerScore = 0
	default:
		examPaperQuestionCustomerAnswer.Answer = &customerQuestionAnswer.Content
		examPaperQuestionCustomerAnswer.CustomerScore = 0
	}
}
func (examService *ExamService) ExamPaperAnswerFromVM(examPaperSubmitVM systemMod.ExamPaperSubmitVM, examPaper systemMod.ExamPaper, examPaperQuestionCustomerAnswers []systemMod.ExamPaperQuestionCustomerAnswer, user *systemMod.SysExamUser, now time.Time) systemMod.ExamPaperAnswer {
	var systemScore int
	var questionCorrect int

	for _, answer := range examPaperQuestionCustomerAnswers {
		systemScore += answer.CustomerScore
		if answer.CustomerScore == answer.QuestionScore {
			questionCorrect++
		}
	}

	examPaperAnswer := systemMod.ExamPaperAnswer{
		PaperName:       examPaper.Name,
		DoTime:          examPaperSubmitVM.DoTime,
		ExamPaperID:     examPaper.ID,
		CreateUser:      2,
		CreateTime:      now,
		SubjectID:       examPaper.SubjectID,
		QuestionCount:   examPaper.QuestionCount,
		PaperScore:      examPaper.Score,
		PaperType:       examPaper.PaperType,
		SystemScore:     systemScore,
		UserScore:       systemScore,
		TaskExamID:      examPaper.TaskExamID,
		QuestionCorrect: questionCorrect,
	}

	needJudge := false
	for _, d := range examPaperQuestionCustomerAnswers {
		if ExamServiceApp.QuestionTypeEnumNeedSaveTextContent(d.QuestionType) {
			needJudge = true
			break
		}
	}

	if needJudge {
		examPaperAnswer.Status = global.ExamPaperAnswerStatusEnumWaitJudge
	} else {
		examPaperAnswer.Status = global.ExamPaperAnswerStatusEnumComplete
	}
	global.GES_LOG.Info("exam3", zap.Any("ExamPaperAnswerFromVM", examPaperAnswer))
	return examPaperAnswer
}

func (examService *ExamService) QuestionTypeEnumNeedSaveTextContent(questionType int) bool {
	// Implement logic based on your enum
	return false
}

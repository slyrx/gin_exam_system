package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/model/common/response"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
	"go.uber.org/zap"
)

type ExamApi struct{}

func (e *ExamApi) AnswerSubmit(c *gin.Context) {
	var examPaperSubmitVM systemMod.ExamPaperSubmitVM
	err := c.ShouldBindJSON(&examPaperSubmitVM)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	global.GES_LOG.Info("exam", zap.Any("", examPaperSubmitVM))

	// 获取当前用户
	user := utils.GetCurrentUser(c)
	// 计算试卷答案信息
	examPaperAnswerInfo, err := examService.CalculateExamPaperAnswer(user, examPaperSubmitVM)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 创建用户事件日志对象
	userEventLog := systemMod.UserEventLog{
		UserID:     1,
		UserName:   "student",
		RealName:   "学生",
		Content:    "student 提交试卷：" + examPaperAnswerInfo.ExamPaper.Name + " 得分：" + utils.ScoreToVM(examPaperAnswerInfo.ExamPaperAnswer.UserScore) + " 耗时：" + utils.SecondToVM(examPaperAnswerInfo.ExamPaperAnswer.DoTime),
		CreateTime: time.Now(),
	}
	// 保存用户事件日志
	examService.CreateUserEventLog(userEventLog)
	response.OkWithMessageExam(utils.ScoreToVM(examPaperAnswerInfo.ExamPaperAnswer.UserScore), c)
}

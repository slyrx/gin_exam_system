package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/model/common/response"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/model/system/request"
	"github.com/slyrx/gin_exam_system/server/others/utils"
)

type ExamApi struct{}

func (h *ExamApi) AnswerSubmit(c *gin.Context) {
	var req request.AnswerSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := examService.GetUserInfo(utils.GetCurrentUser(c).UserName).ID // 假设这个函数已经实现

	totalScore, doTime, paperName, err := examService.SubmitAnswer(req, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 创建用户事件日志对象
	userEventLog := systemMod.UserEventLog{
		UserID:     1,
		UserName:   "student",
		RealName:   "学生",
		Content:    "student 提交试卷：" + paperName + " 得分：" + utils.ScoreToVM(totalScore) + " 耗时：" + utils.SecondToVM(doTime),
		CreateTime: time.Now(),
	}
	// 保存用户事件日志
	examService.CreateUserEventLog(userEventLog)
	response.OkWithMessageExam(utils.ScoreToVM(totalScore), c)
}

package system

import (
	"errors"
	"strconv"

	"github.com/slyrx/gin_exam_system/server/model/common/response"
	"github.com/slyrx/gin_exam_system/server/model/system/request"
	sysModRes "github.com/slyrx/gin_exam_system/server/model/system/response"

	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/others/utils"
)

type ExamPaperApi struct{}

func (e *ExamPaperApi) SetExamPaper(c *gin.Context) {
	var req request.CreateExamPaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证 titleItems 长度
	if len(req.TitleItems) == 0 {
		response.FailWithMessage("titleItems 长度不能为 0", c)
		return
	}

	// 解析 SuggestTime
	suggestTime, err := e.parseSuggestTime(req.SuggestTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.SuggestTime = suggestTime
	userID := examService.GetUserInfo(utils.GetCurrentUser(c).UserName).ID // 假设这个函数已经实现

	// 调用 service 层处理业务逻辑
	paperID, err := examPaperService.CreateExamPaper(userID, req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessageExamInterface(sysModRes.CreateExamPaperResponse{
		ID:            paperID,
		Level:         req.Level,
		SubjectID:     req.SubjectID,
		PaperType:     req.PaperType,
		Name:          req.Name,
		SuggestTime:   suggestTime,
		LimitDateTime: req.LimitDateTime,
		TitleItems:    req.TitleItems,
		Score:         "10", // 假设总分为10，实际应计算总分
	}, c)
}

func (e *ExamPaperApi) parseSuggestTime(suggestTime interface{}) (int, error) {
	switch v := suggestTime.(type) {
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	default:
		return 0, errors.New("无效的 SuggestTime 类型")
	}
}

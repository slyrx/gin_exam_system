package system

import (
	"errors"
	"net/http"
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

func (e *ExamPaperApi) CreateErrorQuestionPaper(c *gin.Context) {
	var req request.CreateErrorQuestionPaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := examService.GetUserInfo(utils.GetCurrentUser(c).UserName).ID // 假设这个函数已经实现
	paperID, err := examPaperService.CreateErrorQuestionPaper(req.SubjectID, req.GradeLevel, userID)
	if err != nil {
		response.FailWithMessage("Failed to create exam paper", c)
		return
	}

	response.OkWithMessageExamInterface(sysModRes.CreateErrorQuestionPaperResponse{PaperID: paperID}, c)
}

func (e *ExamPaperApi) CreateErrorQuestionPaperByUser(c *gin.Context) {
	var req request.CreateErrorQuestionPaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := examService.GetUserInfo(utils.GetCurrentUser(c).UserName).ID // 假设这个函数已经实现
	paperID, err := examPaperService.CreateErrorQuestionPaperByUser(req.SubjectID, req.GradeLevel, userID, req.UserID)
	if err != nil {
		response.FailWithMessage("Failed to create exam paper", c)
		return
	}

	response.OkWithMessageExamInterface(sysModRes.CreateErrorQuestionPaperResponse{PaperID: paperID}, c)
}

func (h *ExamPaperApi) AssignPaperVisibility(c *gin.Context) {
	var req request.AssignPaperVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID (假设我们已经在中间件中设置了用户信息)
	userID := examService.GetUserInfo(utils.GetCurrentUser(c).UserName).ID

	if err := examPaperService.AssignPaperVisibility(req, userID); err != nil {
		response.FailWithMessage("Failed to assign paper visibility", c)
		return
	}

	response.OkWithMessageExam("成功", c)
}
